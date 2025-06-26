package dao

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/supuwoerc/weaver/pkg/database"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestNewBasicDao(t *testing.T) {
	t.Run("successful creation with valid DB", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer func() {
			_ = db.Close()
		}()
		defer func() {
			assert.NoError(t, mock.ExpectationsWereMet())
		}()
		gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}))
		require.NoError(t, err)
		basicDAO := NewBasicDao(gormDB)
		assert.NotNil(t, basicDAO)
		assert.Same(t, gormDB, basicDAO.DB)
	})

	t.Run("creation with nil DB", func(t *testing.T) {
		basicDAO := NewBasicDao(nil)
		assert.NotNil(t, basicDAO)
		assert.Nil(t, basicDAO.DB)
	})
}

type BasicDAOSuite struct {
	basicDAO *BasicDAO
	mock     sqlmock.Sqlmock
	db       *sql.DB
	gormDB   *gorm.DB
	suite.Suite
}

func mockCountArgs(count int) []driver.Value {
	arguments := make([]driver.Value, count)
	for i := range arguments {
		arguments[i] = sqlmock.AnyArg()
	}
	return arguments
}

func (s *BasicDAOSuite) SetupSuite() {
	t := s.T()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}))
	require.NoError(t, err)
	s.basicDAO = NewBasicDao(gormDB)
	s.mock = mock
	s.db = db
	s.gormDB = gormDB
}

func (s *BasicDAOSuite) TearDownSuite() {
	_ = s.db.Close()
}

func TestBasicDAOSuite(t *testing.T) {
	suite.Run(t, new(BasicDAOSuite))
}

func (s *BasicDAOSuite) TestBasicDAO_Datasource() {
	t := s.T()
	s.Run("with nil context", func() {
		var ctx context.Context = nil
		result := s.basicDAO.Datasource(ctx)
		assert.Equal(t, s.gormDB, result)

	})

	s.Run("with context without transaction manager", func() {
		ctx := context.Background()
		result := s.basicDAO.Datasource(ctx)
		assert.NotNil(t, result)
		assert.NotEqual(t, s.gormDB, result) // 应该是一个新的 DB 实例

	})

	s.Run("with context containing transaction manager", func() {
		manager := &database.TransactionManager{
			DB: s.gormDB,
		}
		injectManagerCtx := database.InjectManager(context.Background(), manager)
		result := s.basicDAO.Datasource(injectManagerCtx)
		assert.NotNil(t, result)
		assert.Equal(t, manager.DB, s.gormDB)
		assert.NotEqual(t, s.gormDB, result)
		assert.NotEqual(t, manager.DB, result)
	})
}

func (s *BasicDAOSuite) TestBasicDAO_Preload() {
	t := s.T()
	s.Run("preload with single field", func() {
		defer func() {
			assert.NoError(t, s.mock.ExpectationsWereMet())
		}()
		preloadFunc := s.basicDAO.Preload("users")
		assert.IsType(t, (func(d *gorm.DB) *gorm.DB)(nil), preloadFunc)
		result := preloadFunc(s.gormDB)
		assert.NotNil(t, result)
	})

	s.Run("preload with field and arguments", func() {
		defer func() {
			assert.NoError(t, s.mock.ExpectationsWereMet())
		}()
		preloadFunc := s.basicDAO.Preload("users", "status = ?", "active")
		assert.IsType(t, (func(d *gorm.DB) *gorm.DB)(nil), preloadFunc)
		result := preloadFunc(s.gormDB)
		assert.NotNil(t, result)
	})
}
