package service

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/supuwoerc/weaver/conf"
	"github.com/supuwoerc/weaver/initialize"
	"github.com/supuwoerc/weaver/pkg/database"
	"github.com/supuwoerc/weaver/pkg/logger"
	"github.com/supuwoerc/weaver/pkg/utils"
	"go.uber.org/zap/zaptest"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestNewBasicService(t *testing.T) {
	testCases := []struct {
		name        string
		logger      *logger.Logger
		db          *gorm.DB
		locksmith   *utils.RedisLocksmith
		conf        *conf.Config
		emailClient *initialize.EmailClient
		wantNil     bool
	}{
		{
			name:        "all dependencies provided",
			logger:      &logger.Logger{},
			db:          &gorm.DB{},
			locksmith:   &utils.RedisLocksmith{},
			conf:        &conf.Config{},
			emailClient: &initialize.EmailClient{},
			wantNil:     false,
		},
		{
			name:        "nil logger",
			logger:      nil,
			db:          &gorm.DB{},
			locksmith:   &utils.RedisLocksmith{},
			conf:        &conf.Config{},
			emailClient: &initialize.EmailClient{},
			wantNil:     false,
		},
		{
			name:        "nil DB",
			logger:      &logger.Logger{},
			db:          nil,
			locksmith:   &utils.RedisLocksmith{},
			conf:        &conf.Config{},
			emailClient: &initialize.EmailClient{},
			wantNil:     false,
		},
		{
			name:        "nil locksmith",
			logger:      &logger.Logger{},
			db:          &gorm.DB{},
			locksmith:   nil,
			conf:        &conf.Config{},
			emailClient: &initialize.EmailClient{},
			wantNil:     false,
		},
		{
			name:        "nil config",
			logger:      &logger.Logger{},
			db:          &gorm.DB{},
			locksmith:   &utils.RedisLocksmith{},
			conf:        nil,
			emailClient: &initialize.EmailClient{},
			wantNil:     false,
		},
		{
			name:        "nil email client",
			logger:      &logger.Logger{},
			db:          &gorm.DB{},
			locksmith:   &utils.RedisLocksmith{},
			conf:        &conf.Config{},
			emailClient: nil,
			wantNil:     false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := NewBasicService(
				tc.logger,
				tc.db,
				tc.locksmith,
				tc.conf,
				tc.emailClient,
			)
			if tc.wantNil {
				assert.Nil(t, service, "service should be nil")
			} else {
				assert.NotNil(t, service, "service should not be nil")
				assert.Equal(t, tc.logger, service.Logger, "logger not properly set")
				assert.Equal(t, tc.db, service.DB, "DB not properly set")
				assert.Equal(t, tc.locksmith, service.Locksmith, "locksmith not properly set")
				assert.Equal(t, tc.conf, service.Conf, "config not properly set")
				assert.Equal(t, tc.emailClient, service.EmailClient, "email client not properly set")
			}
		})
	}
}

type BasicServiceSuite struct {
	suite.Suite
	sqlDB   *sql.DB
	db      *gorm.DB
	mock    sqlmock.Sqlmock
	logger  *logger.Logger
	service *BasicService
}

func (s *BasicServiceSuite) SetupSuite() {
	t := s.T()
	mockLogger := logger.NewLogger(zaptest.NewLogger(s.T()).Sugar())
	// 把匹配器设置成相等匹配器，不设置默认使用正则匹配
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	mockGormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		DisableAutomaticPing: true,
	})
	require.NoError(t, err)
	s.sqlDB = db
	s.db = mockGormDB
	s.mock = mock
	s.logger = mockLogger
	s.service = NewBasicService(mockLogger, mockGormDB, nil, nil, nil)
}

func (s *BasicServiceSuite) TearDownSuite() {
	_ = s.sqlDB.Close()
}

func TestBasicServiceSuite(t *testing.T) {
	suite.Run(t, new(BasicServiceSuite))
}

func (s *BasicServiceSuite) TestBasicService_TransactionErrorAndPanic() {
	var err error
	t := s.T()
	s.Run("Simple Transaction", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectCommit()
		defer func() {
			assert.NoError(t, s.mock.ExpectationsWereMet())
		}()
		err = s.service.Transaction(context.Background(), false, func(ctx context.Context) error {
			return nil
		})
		assert.NoError(t, err)
	})
	s.Run("Nested Join Transaction", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectCommit()
		defer func() {
			assert.NoError(t, s.mock.ExpectationsWereMet())
		}()
		err = s.service.Transaction(context.Background(), false, func(ctx context.Context) error {
			return s.service.Transaction(ctx, true, func(ctx context.Context) error {
				return nil
			})
		})
		assert.NoError(t, err)
	})
	s.Run("Rollback On Error", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectRollback()
		defer func() {
			assert.NoError(t, s.mock.ExpectationsWereMet())
		}()
		e := fmt.Errorf("test error")
		err = s.service.Transaction(context.Background(), false, func(ctx context.Context) error {
			return e
		})
		assert.ErrorContains(t, err, e.Error())
		assert.Equal(t, err, e)
	})
	s.Run("Panic Recovery", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectRollback()
		defer func() {
			assert.NoError(t, s.mock.ExpectationsWereMet())
		}()
		err = s.service.Transaction(context.Background(), false, func(ctx context.Context) error {
			panic("test panic")
		})
		assert.Error(t, err)
		assert.ErrorContains(t, err, "test panic")
	})
	s.Run("Transaction With Options", func() {
		opts := &sql.TxOptions{
			Isolation: sql.LevelSerializable,
			ReadOnly:  true,
		}
		s.mock.ExpectBegin()
		s.mock.ExpectCommit()
		defer func() {
			assert.NoError(t, s.mock.ExpectationsWereMet())
		}()
		err = s.service.Transaction(context.Background(), false, func(ctx context.Context) error {
			return nil
		}, opts)
		assert.NoError(t, err)
	})
	s.Run("Context Propagation", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectCommit()
		defer func() {
			assert.NoError(t, s.mock.ExpectationsWereMet())
		}()
		var capturedCtx context.Context
		err = s.service.Transaction(context.Background(), false, func(ctx context.Context) error {
			capturedCtx = ctx
			return nil
		})
		assert.NoError(t, err)
		manager := database.LoadManager(capturedCtx)
		assert.NotNil(t, manager)
	})
	s.Run("join=true without existing transaction", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectCommit()
		defer func() {
			assert.NoError(t, s.mock.ExpectationsWereMet())
		}()
		ctx := context.Background()
		var capturedManager *database.TransactionManager
		err = s.service.Transaction(ctx, true, func(newCtx context.Context) error {
			// 在事务函数内部捕获 TransactionManager
			capturedManager = database.LoadManager(newCtx)
			require.NotNil(t, capturedManager)
			// 验证这是一个新的事务（isStarter = true 的效果）
			require.False(t, capturedManager.AlreadyCommittedOrRolledBack)
			return nil
		})
		assert.NoError(t, err)
		assert.NotNil(t, capturedManager)
		assert.True(t, capturedManager.AlreadyCommittedOrRolledBack)
	})
	s.Run("Begin error", func() {
		// 模拟 Begin 操作失败
		expectedErr := fmt.Errorf("begin transaction error")
		s.mock.ExpectBegin().WillReturnError(expectedErr)
		defer func() {
			assert.NoError(t, s.mock.ExpectationsWereMet())
		}()
		err = s.service.Transaction(context.Background(), false, func(ctx context.Context) error {
			t.Fatal("transaction function should not be called")
			return nil
		})
		assert.Equal(t, err, expectedErr)
	})
	s.Run("Rollback error after execution error", func() {
		// 模拟 Rollback 操作失败
		expectedErr := fmt.Errorf("rollback fail")
		s.mock.ExpectBegin()
		s.mock.ExpectRollback().WillReturnError(expectedErr)
		defer func() {
			assert.NoError(t, s.mock.ExpectationsWereMet())
		}()
		err = s.service.Transaction(context.Background(), false, func(ctx context.Context) error {
			// 返回一个错误触发回滚
			return fmt.Errorf("inner execution error")
		})
		assert.ErrorContains(t, err, expectedErr.Error())
	})
	s.Run("Commit error", func() {
		// 模拟 Commit 操作失败
		expectedErr := fmt.Errorf("commit fail")
		s.mock.ExpectBegin()
		s.mock.ExpectCommit().WillReturnError(expectedErr)
		defer func() {
			assert.NoError(t, s.mock.ExpectationsWereMet())
		}()
		err = s.service.Transaction(context.Background(), false, func(ctx context.Context) error {
			return nil
		})
		assert.Equal(t, err, expectedErr)
	})
}

type TestUser struct {
	ID   uint
	Name string
}

func (s *BasicServiceSuite) TestBasicService_Transaction() {
	var err error
	t := s.T()
	// 发生错误回滚
	s.Run("rollback on execution error", func() {
		u := TestUser{
			Name: "test name",
		}
		ctx := context.Background()
		insertRaw := `INSERT INTO test_users (name) VALUES (?)`
		queryRaw := `SELECT COUNT(*) FROM test_users`
		var count int64
		s.mock.ExpectExec(insertRaw).WithArgs(u.Name).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectQuery(queryRaw).WillReturnRows(sqlmock.NewRows([]string{"num"}).AddRow(1))
		s.mock.ExpectBegin()
		s.mock.ExpectExec(insertRaw).
			WithArgs(u.Name).WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectRollback()
		s.mock.ExpectQuery(queryRaw).WillReturnRows(sqlmock.NewRows([]string{"num"}).AddRow(1))
		defer func() {
			assert.NoError(t, s.mock.ExpectationsWereMet())
		}()
		beforeTx := s.db.WithContext(ctx)
		err = beforeTx.Exec(insertRaw, u.Name).Error
		require.NoError(t, err)
		err = beforeTx.Raw(queryRaw).Scan(&count).Error
		require.NoError(t, err)
		assert.Equal(t, count, int64(1))
		count = 0 // reset count
		err = s.service.Transaction(ctx, false, func(ctx context.Context) error {
			tx := s.db.WithContext(ctx)
			e := tx.Exec(insertRaw, u.Name).Error
			require.NoError(t, e)
			return fmt.Errorf("force fail")
		})
		assert.Error(t, err)
		assert.ErrorContains(t, err, "force fail")
		err = beforeTx.Raw(queryRaw).Scan(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, count, int64(1))
	})

	// 发生panic回滚
	s.Run("rollback on panic", func() {
		u := TestUser{
			Name: "test name",
		}
		ctx := context.Background()
		insertRaw := `INSERT INTO test_users (name) VALUES (?)`
		queryRaw := `SELECT COUNT(*) FROM test_users`
		var count int64
		s.mock.ExpectExec(insertRaw).WithArgs(u.Name).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectQuery(queryRaw).WillReturnRows(sqlmock.NewRows([]string{"num"}).AddRow(1))
		s.mock.ExpectBegin()
		s.mock.ExpectExec(insertRaw).WithArgs(u.Name).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectRollback()
		s.mock.ExpectQuery(queryRaw).WillReturnRows(sqlmock.NewRows([]string{"num"}).AddRow(1))
		defer func() {
			assert.NoError(t, s.mock.ExpectationsWereMet())
		}()
		beforeTx := s.db.WithContext(ctx)
		err = beforeTx.Exec(insertRaw, u.Name).Error
		require.NoError(t, err)
		err = beforeTx.Raw(queryRaw).Scan(&count).Error
		require.NoError(t, err)
		assert.Equal(t, count, int64(1))
		count = 0 // reset count
		err = s.service.Transaction(ctx, false, func(ctx context.Context) error {
			tx := s.db.WithContext(ctx)
			e := tx.Exec(insertRaw, u.Name).Error
			require.NoError(t, e)
			panic("transaction with panic")
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction with panic")
		err = beforeTx.Raw(queryRaw).Scan(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, count, int64(1))
	})

	// 嵌套事务的回滚
	s.Run("nested transaction rollback", func() {
		u := TestUser{
			Name: "test name",
		}
		u2 := TestUser{
			Name: "test name 2",
		}
		ctx := context.Background()
		insertRaw := `INSERT INTO test_users (name) VALUES (?)`
		queryRaw := `SELECT COUNT(*) FROM test_users`
		var count int64
		s.mock.ExpectExec(insertRaw).WithArgs(u.Name).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectQuery(queryRaw).WillReturnRows(sqlmock.NewRows([]string{"num"}).AddRow(1))
		s.mock.ExpectBegin()
		s.mock.ExpectExec(insertRaw).WithArgs(u.Name).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectExec(insertRaw).WithArgs(u2.Name).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectRollback()
		s.mock.ExpectQuery(queryRaw).WillReturnRows(sqlmock.NewRows([]string{"num"}).AddRow(1))
		defer func() {
			assert.NoError(t, s.mock.ExpectationsWereMet())
		}()
		beforeTx := s.db.WithContext(ctx)
		err = beforeTx.Exec(insertRaw, u.Name).Error
		require.NoError(t, err)
		err = beforeTx.Raw(queryRaw).Scan(&count).Error
		require.NoError(t, err)
		assert.Equal(t, count, int64(1))
		count = 0 // reset count
		err = s.service.Transaction(ctx, false, func(ctx context.Context) error {
			tx := s.db.WithContext(ctx)
			e := tx.Exec(insertRaw, u.Name).Error
			require.NoError(t, e)
			return s.service.Transaction(ctx, true, func(ctx context.Context) error {
				e = tx.Exec(insertRaw, u2.Name).Error
				require.NoError(t, e)
				return fmt.Errorf("nested transaction error")
			})
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nested transaction error")
		err = beforeTx.Raw(queryRaw).Scan(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, count, int64(1))
	})

	// 测试回滚失败的情况
	s.Run("rollback failure", func() {
		u := TestUser{
			Name: "test name",
		}
		ctx := context.Background()
		insertRaw := `INSERT INTO test_users (name) VALUES (?)`
		queryRaw := `SELECT COUNT(*) FROM test_users`
		var count int64
		s.mock.ExpectExec(insertRaw).WithArgs(u.Name).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectQuery(queryRaw).WillReturnRows(sqlmock.NewRows([]string{"num"}).AddRow(1))
		s.mock.ExpectBegin()
		s.mock.ExpectExec(insertRaw).WithArgs(u.Name).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectRollback().WillReturnError(fmt.Errorf("rollback fail"))
		s.mock.ExpectQuery(queryRaw).WillReturnRows(sqlmock.NewRows([]string{"num"}).AddRow(2))
		defer func() {
			assert.NoError(t, s.mock.ExpectationsWereMet())
		}()
		beforeTx := s.db.WithContext(ctx)
		err = beforeTx.Exec(insertRaw, u.Name).Error
		require.NoError(t, err)
		err = beforeTx.Raw(queryRaw).Scan(&count).Error
		require.NoError(t, err)
		assert.Equal(t, count, int64(1))
		count = 0 // reset count
		err = s.service.Transaction(ctx, false, func(ctx context.Context) error {
			tx := s.db.WithContext(ctx)
			e := tx.Exec(insertRaw, u.Name).Error
			require.NoError(t, e)
			return fmt.Errorf("force rollback")
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rollback fail")
		assert.Contains(t, err.Error(), "force rollback")
		err = beforeTx.Raw(queryRaw).Scan(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, count, int64(2))
	})

	// 测试部分提交后的回滚
	s.Run("rollback after partial commit", func() {
		u := TestUser{
			Name: "test name",
		}
		updateName := "update name"
		ctx := context.Background()
		insertRaw := `INSERT INTO test_users (name) VALUES (?)`
		updateRaw := `UPDATE test_users set name = ?`
		queryRaw := `SELECT COUNT(*) FROM test_users`
		var count int64
		s.mock.ExpectExec(insertRaw).WithArgs(u.Name).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectQuery(queryRaw).WillReturnRows(sqlmock.NewRows([]string{"num"}).AddRow(1))
		s.mock.ExpectBegin()
		s.mock.ExpectExec(insertRaw).WithArgs(u.Name).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectExec(updateRaw).WithArgs(updateName).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectRollback()
		s.mock.ExpectQuery(queryRaw).WillReturnRows(sqlmock.NewRows([]string{"num"}).AddRow(1))
		defer func() {
			assert.NoError(t, s.mock.ExpectationsWereMet())
		}()
		beforeTx := s.db.WithContext(ctx)
		err = beforeTx.Exec(insertRaw, u.Name).Error
		require.NoError(t, err)
		err = beforeTx.Raw(queryRaw).Scan(&count).Error
		require.NoError(t, err)
		assert.Equal(t, count, int64(1))
		count = 0 // reset count
		err = s.service.Transaction(ctx, false, func(ctx context.Context) error {
			tx := s.db.WithContext(ctx)
			e := tx.Exec(insertRaw, u.Name).Error
			require.NoError(t, e)
			e = tx.Exec(updateRaw, updateName).Error
			require.NoError(t, e)
			return fmt.Errorf("force rollback")
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "force rollback")
		err = beforeTx.Raw(queryRaw).Scan(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, count, int64(1))
	})
}
