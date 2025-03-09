package dao

import (
	"context"
	"errors"
	"gin-web/models"
	"gin-web/pkg/constant"
	"gin-web/pkg/database"
	"gin-web/pkg/response"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type UserDAO struct {
	*BasicDAO
}

func NewUserDAO(basicDAO *BasicDAO) *UserDAO {
	return &UserDAO{
		BasicDAO: basicDAO,
	}
}

func (u *UserDAO) Create(ctx context.Context, user *models.User) error {
	err := u.Datasource(ctx).Create(user).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return response.UserCreateDuplicateEmail
	}
	return err
}

func (u *UserDAO) GetByEmail(ctx context.Context, email string, needRoles, needPermissions, needDepts bool) (*models.User, error) {
	var user models.User
	query := u.Datasource(ctx).Model(&models.User{})
	if needRoles {
		query = query.Preload("Roles")
		if needPermissions {
			query.Preload("Roles.Permissions")
		}
	}
	if needDepts {
		query = query.Preload("Departments")
	}
	err := query.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.UserNotExist
		}
		return nil, err
	}
	return &user, nil
}

func (u *UserDAO) GetById(ctx context.Context, uid uint, needRoles, needPermissions, needDepts bool) (*models.User, error) {
	var user models.User
	query := u.Datasource(ctx).Model(&models.User{})
	if needRoles {
		query.Preload("Roles")
		if needPermissions {
			query.Preload("Roles.Permissions")
		}
	}
	if needDepts {
		query = query.Preload("Departments")
	}
	err := query.Where("id = ?", uid).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.UserNotExist
		}
		return nil, err
	}
	return &user, nil
}

func (u *UserDAO) GetByIds(ctx context.Context, ids []uint, needRoles, needPermissions, needDepts bool) ([]*models.User, error) {
	var users []*models.User
	query := u.Datasource(ctx).Model(&models.User{})
	if needRoles {
		query.Preload("Roles")
		if needPermissions {
			query.Preload("Roles.Permissions")
		}
	}
	if needDepts {
		query = query.Preload("Departments")
	}
	err := query.Where("id in (?)", ids).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserDAO) GetList(ctx context.Context, keyword string, limit, offset int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64
	query := u.Datasource(ctx).Model(&models.User{}).Order("updated_at desc,id desc")
	if keyword != "" {
		keyword = database.FuzzKeyword(keyword)
		query = query.Where("nickname like ? or email like ?", keyword, keyword)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Preload("Roles").Preload("Departments").Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (u *UserDAO) GetAll(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	err := u.Datasource(ctx).Model(&models.User{}).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserDAO) UpdateAccountStatus(ctx context.Context, id uint, status constant.UserStatus) error {
	return u.Datasource(ctx).Model(&models.User{}).Where("id = ?", id).Update("status", status).Error
}
