package dao

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	sqlDriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/supuwoerc/weaver/models"
	"github.com/supuwoerc/weaver/pkg/database"
	"github.com/supuwoerc/weaver/pkg/response"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestNewDepartmentDAO(t *testing.T) {
	t.Run("successful creation with valid BasicDAO", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer func() {
			_ = db.Close()
		}()
		defer func() {
			assert.NoError(t, mock.ExpectationsWereMet())
		}()
		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}))
		require.NoError(t, err)
		basicDAO := NewBasicDao(gormDB)
		deptDAO := NewDepartmentDAO(basicDAO)
		assert.NotNil(t, basicDAO)
		assert.NotNil(t, deptDAO)
		assert.Equal(t, deptDAO.BasicDAO, basicDAO)
	})

	t.Run("creation with nil BasicDAO", func(t *testing.T) {
		deptDAO := NewDepartmentDAO(nil)
		assert.NotNil(t, deptDAO)
		assert.Nil(t, deptDAO.BasicDAO)
	})

	t.Run("creation with BasicDAO having nil DB", func(t *testing.T) {
		basicDAO := &BasicDAO{DB: nil}
		deptDAO := NewDepartmentDAO(basicDAO)
		assert.NotNil(t, deptDAO)
		assert.Equal(t, basicDAO, deptDAO.BasicDAO)
		assert.Nil(t, deptDAO.DB)
	})

	t.Run("verify method inheritance", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer func() {
			_ = db.Close()
		}()
		defer func() {
			assert.NoError(t, mock.ExpectationsWereMet())
		}()
		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		require.NoError(t, err)
		basicDAO := &BasicDAO{DB: gormDB}
		deptDAO := NewDepartmentDAO(basicDAO)
		ctx := context.Background()
		datasource := deptDAO.Datasource(ctx)
		assert.NotNil(t, datasource)
		assert.IsType(t, &gorm.DB{}, datasource)
	})
}

type DepartmentDAOSuite struct {
	deptDAO *DepartmentDAO
	mock    sqlmock.Sqlmock
	db      *sql.DB
	gormDB  *gorm.DB
	suite.Suite
}

func TestDepartmentDAOSuite(t *testing.T) {
	suite.Run(t, new(DepartmentDAOSuite))
}

func (s *DepartmentDAOSuite) SetupSuite() {
	t := s.T()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}))
	require.NoError(t, err)
	s.deptDAO = NewDepartmentDAO(NewBasicDao(gormDB))
	s.mock = mock
	s.db = db
	s.gormDB = gormDB
}

func (s *DepartmentDAOSuite) TearDownSuite() {
	_ = s.db.Close()
}

func (s *DepartmentDAOSuite) TestDepartmentDAO_Create() {
	t := s.T()
	s.Run("successful create department", func() {
		dept := &models.Department{
			Name:      "IT部门",
			CreatorId: 1,
			UpdaterId: 1,
			BasicModel: database.BasicModel{
				CreatedAt: database.UpsertTime(time.Now()),
				UpdatedAt: database.UpsertTime(time.Now()),
			},
		}
		defer func() {
			assert.NoError(t, s.mock.ExpectationsWereMet())
		}()
		s.mock.ExpectBegin()
		s.mock.ExpectExec("INSERT INTO `departments`").
			WithArgs(mockCountArgs(8)...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()
		err := s.deptDAO.Create(context.Background(), dept)
		assert.NoError(t, err)
	})

	s.Run("duplicate name error", func() {
		dept := &models.Department{
			Name:      "IT部门",
			CreatorId: 1,
			UpdaterId: 1,
		}
		defer func() {
			assert.NoError(t, s.mock.ExpectationsWereMet())
		}()
		s.mock.ExpectBegin()
		s.mock.ExpectExec("INSERT INTO `departments`").
			WithArgs(mockCountArgs(8)...).
			WillReturnError(&sqlDriver.MySQLError{
				Number:  1062,
				Message: "Test Error: Duplicate entry 'IT部门' for uniq index",
			})
		s.mock.ExpectRollback()
		err := s.deptDAO.Create(context.Background(), dept)
		assert.Error(t, err)
		assert.Equal(t, response.RoleCreateDuplicateName, err)
	})

	s.Run("other db error", func() {
		dept := &models.Department{
			Name:      "IT部门",
			CreatorId: 1,
			UpdaterId: 1,
		}
		defer func() {
			assert.NoError(t, s.mock.ExpectationsWereMet())
		}()
		s.mock.ExpectBegin()
		s.mock.ExpectExec("INSERT INTO `departments`").
			WithArgs(mockCountArgs(8)...).
			WillReturnError(&sqlDriver.MySQLError{
				Number:  999,
				Message: "Test Error: DB error",
			})
		s.mock.ExpectRollback()
		err := s.deptDAO.Create(context.Background(), dept)
		assert.Error(t, err)
		assert.NotEqual(t, response.RoleCreateDuplicateName, err)
	})
}

func (s *DepartmentDAOSuite) TestDepartmentDAO_GetByName() {
	t := s.T()
	s.Run("successful get by name", func() {
		expectedDept := &models.Department{
			BasicModel: database.BasicModel{
				ID: 1998,
			},
			Name:      "IT部门",
			CreatorId: 1,
			UpdaterId: 1,
		}
		defer func() {
			assert.NoError(t, s.mock.ExpectationsWereMet())
		}()
		queryRaw := "SELECT \\* FROM `departments` WHERE name = \\? AND `departments`.`deleted_at` = \\? ORDER BY `departments`.`id` LIMIT \\?"
		s.mock.ExpectQuery(queryRaw).
			WithArgs("IT部门", 0, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "creator_id", "updater_id"}).
				AddRow(1998, "IT部门", 1, 1))
		dept, err := s.deptDAO.GetByName(context.Background(), "IT部门")
		assert.NoError(t, err)
		assert.NotNil(t, dept)
		assert.Equal(t, expectedDept.Name, dept.Name)
		assert.Equal(t, expectedDept.ID, dept.ID)
		assert.Equal(t, expectedDept, dept)
	})

	s.Run("department not found", func() {
		defer func() {
			assert.NoError(t, s.mock.ExpectationsWereMet())
		}()
		queryRaw := "SELECT \\* FROM `departments` WHERE name = \\? AND `departments`.`deleted_at` = \\? ORDER BY `departments`.`id` LIMIT \\?"
		s.mock.ExpectQuery(queryRaw).
			WithArgs("不存在的部门", 0, 1).
			WillReturnError(gorm.ErrRecordNotFound)
		dept, err := s.deptDAO.GetByName(context.Background(), "不存在的部门")
		assert.Error(t, err)
		assert.Nil(t, dept)
		assert.Equal(t, response.DeptNotExist, err)
	})
}
