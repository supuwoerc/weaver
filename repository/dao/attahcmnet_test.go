package dao

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/supuwoerc/weaver/models"
	"github.com/supuwoerc/weaver/pkg/database"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestNewAttachmentDAO(t *testing.T) {
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
		deptDAO := NewAttachmentDAO(basicDAO)
		assert.NotNil(t, basicDAO)
		assert.NotNil(t, deptDAO)
		assert.Equal(t, deptDAO.BasicDAO, basicDAO)
	})

	t.Run("creation with nil BasicDAO", func(t *testing.T) {
		deptDAO := NewAttachmentDAO(nil)
		assert.NotNil(t, deptDAO)
		assert.Nil(t, deptDAO.BasicDAO)
	})

	t.Run("creation with BasicDAO having nil DB", func(t *testing.T) {
		basicDAO := &BasicDAO{DB: nil}
		deptDAO := NewAttachmentDAO(basicDAO)
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
		deptDAO := NewAttachmentDAO(basicDAO)
		ctx := context.Background()
		datasource := deptDAO.Datasource(ctx)
		assert.NotNil(t, datasource)
		assert.IsType(t, &gorm.DB{}, datasource)
	})
}

type AttachmentDAOSuite struct {
	attachmentDAO *AttachmentDAO
	mock          sqlmock.Sqlmock
	db            *sql.DB
	gormDB        *gorm.DB
	suite.Suite
}

func (s *AttachmentDAOSuite) SetupSuite() {
	t := s.T()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}))
	require.NoError(t, err)
	s.attachmentDAO = NewAttachmentDAO(NewBasicDao(gormDB))
	s.mock = mock
	s.db = db
	s.gormDB = gormDB
}

func (s *AttachmentDAOSuite) TearDownSuite() {
	_ = s.db.Close()
}

func TestAttachmentDAOSuite(t *testing.T) {
	suite.Run(t, new(AttachmentDAOSuite))
}

func (s *AttachmentDAOSuite) TestAttachmentDAO_Create() {
	t := s.T()
	s.Run("successful creation of single attachment", func() {
		defer func() {
			assert.NoError(t, s.mock.ExpectationsWereMet())
		}()
		attachments := []*models.Attachment{
			{
				Name:      "test.txt",
				Type:      1,
				Size:      1024,
				Hash:      "abc123",
				Path:      "/uploads/test.txt",
				CreatorId: lo.ToPtr[uint](1),
				BasicModel: database.BasicModel{
					CreatedAt: database.UpsertTime(time.Now()),
					UpdatedAt: database.UpsertTime(time.Now()),
				},
			},
		}
		s.mock.ExpectBegin()
		s.mock.ExpectExec("INSERT INTO `attachments`").
			WithArgs(mockCountArgs(9)...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()
		err := s.attachmentDAO.Create(context.Background(), attachments)
		assert.NoError(t, err)
		assert.Equal(t, attachments[0].ID, uint(1))
	})

	s.Run("creation with empty slice", func() {
		defer func() {
			assert.NoError(t, s.mock.ExpectationsWereMet())
		}()
		attachments := make([]*models.Attachment, 0)
		s.mock.ExpectBegin()
		s.mock.ExpectRollback()
		err := s.attachmentDAO.Create(context.Background(), attachments)
		assert.ErrorContains(t, err, "empty slice found")
	})

	s.Run("creation with nil slice", func() {
		defer func() {
			assert.NoError(t, s.mock.ExpectationsWereMet())
		}()
		s.mock.ExpectBegin()
		s.mock.ExpectRollback()
		err := s.attachmentDAO.Create(context.Background(), nil)
		assert.ErrorContains(t, err, "empty slice found")
	})
}
