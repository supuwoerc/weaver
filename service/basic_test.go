package service

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			name:        "nil db",
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
				assert.Equal(t, tc.logger, service.logger, "logger not properly set")
				assert.Equal(t, tc.db, service.db, "db not properly set")
				assert.Equal(t, tc.locksmith, service.locksmith, "locksmith not properly set")
				assert.Equal(t, tc.conf, service.conf, "config not properly set")
				assert.Equal(t, tc.emailClient, service.emailClient, "email client not properly set")
			}
		})
	}
}

func TestBasicService_Transaction(t *testing.T) {
	mockLogger := logger.NewLogger(zaptest.NewLogger(t).Sugar())
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		_ = db.Close()
	}()
	mockGormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		DisableAutomaticPing: true,
	})
	require.NoError(t, err)
	service := NewBasicService(mockLogger, mockGormDB, nil, nil, nil)
	t.Run("Simple Transaction", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectCommit()
		defer func() {
			assert.NoError(t, mock.ExpectationsWereMet())
		}()
		err = service.Transaction(context.Background(), false, func(ctx context.Context) error {
			return nil
		})
		assert.NoError(t, err)
	})
	t.Run("Nested Join Transaction", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectCommit()
		defer func() {
			assert.NoError(t, mock.ExpectationsWereMet())
		}()
		err = service.Transaction(context.Background(), false, func(ctx context.Context) error {
			return service.Transaction(ctx, true, func(ctx context.Context) error {
				return nil
			})
		})
		assert.NoError(t, err)
	})
	t.Run("Rollback On Error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectRollback()
		defer func() {
			assert.NoError(t, mock.ExpectationsWereMet())
		}()
		e := fmt.Errorf("test error")
		err = service.Transaction(context.Background(), false, func(ctx context.Context) error {
			return e
		})
		assert.ErrorContains(t, err, e.Error())
		assert.Equal(t, err, e)
	})
	t.Run("Panic Recovery", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectRollback()
		defer func() {
			assert.NoError(t, mock.ExpectationsWereMet())
		}()
		err = service.Transaction(context.Background(), false, func(ctx context.Context) error {
			panic("test panic")
		})
		assert.Error(t, err)
		assert.ErrorContains(t, err, "test panic")
	})
	t.Run("Transaction With Options", func(t *testing.T) {
		opts := &sql.TxOptions{
			Isolation: sql.LevelSerializable,
			ReadOnly:  true,
		}
		mock.ExpectBegin()
		mock.ExpectCommit()
		defer func() {
			assert.NoError(t, mock.ExpectationsWereMet())
		}()
		err = service.Transaction(context.Background(), false, func(ctx context.Context) error {
			return nil
		}, opts)
		assert.NoError(t, err)
	})
	t.Run("Context Propagation", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectCommit()
		defer func() {
			assert.NoError(t, mock.ExpectationsWereMet())
		}()
		var capturedCtx context.Context
		err = service.Transaction(context.Background(), false, func(ctx context.Context) error {
			capturedCtx = ctx
			return nil
		})
		assert.NoError(t, err)
		manager := database.LoadManager(capturedCtx)
		assert.NotNil(t, manager)
	})
	t.Run("join=true without existing transaction", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectCommit()
		defer func() {
			assert.NoError(t, mock.ExpectationsWereMet())
		}()
		ctx := context.Background()
		var capturedManager *database.TransactionManager
		err = service.Transaction(ctx, true, func(newCtx context.Context) error {
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
	t.Run("Begin error", func(t *testing.T) {
		// 模拟 Begin 操作失败
		expectedErr := fmt.Errorf("begin transaction error")
		mock.ExpectBegin().WillReturnError(expectedErr)
		defer func() {
			assert.NoError(t, mock.ExpectationsWereMet())
		}()
		err = service.Transaction(context.Background(), false, func(ctx context.Context) error {
			t.Fatal("transaction function should not be called")
			return nil
		})
		assert.Equal(t, err, expectedErr)
	})
	t.Run("Rollback error after execution error", func(t *testing.T) {
		// 模拟 Rollback 操作失败
		expectedErr := fmt.Errorf("rollback fail")
		mock.ExpectBegin()
		mock.ExpectRollback().WillReturnError(expectedErr)
		defer func() {
			assert.NoError(t, mock.ExpectationsWereMet())
		}()
		err = service.Transaction(context.Background(), false, func(ctx context.Context) error {
			// 返回一个错误触发回滚
			return fmt.Errorf("inner execution error")
		})
		t.Log(err)
		assert.ErrorContains(t, err, expectedErr.Error())
	})
	t.Run("Commit error", func(t *testing.T) {
		// 模拟 Commit 操作失败
		expectedErr := fmt.Errorf("commit fail")
		mock.ExpectBegin()
		mock.ExpectCommit().WillReturnError(expectedErr)
		defer func() {
			assert.NoError(t, mock.ExpectationsWereMet())
		}()
		err = service.Transaction(context.Background(), false, func(ctx context.Context) error {
			return nil
		})
		assert.Equal(t, err, expectedErr)
	})
}
