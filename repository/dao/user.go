package dao

import (
	"context"
	"errors"
	"gin-web/models"
	"gin-web/pkg/response"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var (
	userDAO     *UserDAO
	userDAOOnce sync.Once
)

type UserDAO struct {
	*BasicDAO
}

func NewUserDAO() *UserDAO {
	userDAOOnce.Do(func() {
		userDAO = &UserDAO{BasicDAO: NewBasicDao()}
	})
	return userDAO
}

func (u *UserDAO) Insert(ctx context.Context, user *models.User) error {
	err := u.Datasource(ctx).Create(user).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return response.UserCreateDuplicateEmail
	}
	return err
}

func (u *UserDAO) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := u.Datasource(ctx).Preload("Roles").Where("email = ?", email).First(&user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return &user, response.UserLoginEmailNotFound
	}
	return &user, err
}

func (u *UserDAO) AssociateRoles(ctx context.Context, uid uint, roles []*models.Role) error {
	return u.Datasource(ctx).Model(&models.User{}).Where("id = ?", uid).Association("Roles").Replace(roles)
}

func (u *UserDAO) FindByUid(ctx context.Context, uid uint, needRoles, needPermissions bool) (*models.User, error) {
	var user models.User
	query := u.Datasource(ctx).Model(&models.User{})
	if needRoles {
		query.Preload("Roles")
		if needPermissions {
			query.Preload("Roles.Permissions")
		}
	}
	err := query.Where("id = ?", uid).First(&user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return &user, response.UserNotExist
	}
	return &user, err
}

func (u *UserDAO) FindRolesByUid(ctx context.Context, uid uint) ([]*models.Role, error) {
	var result []*models.Role
	// TODO:改造用Association查询
	err := u.Datasource(ctx).
		Table("sys_role as r").
		Select("r.id", "r.name").
		Joins("join sys_user_role as ur on r.id = ur.role_id and ur.user_id = ?", uid).
		Scan(&result).Error
	return result, err
}
