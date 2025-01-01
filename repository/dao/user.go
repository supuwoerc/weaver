package dao

import (
	"context"
	"database/sql"
	"errors"
	"gin-web/pkg/response"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type UserDAO struct {
	*BasicDAO
}

type User struct {
	gorm.Model
	Email    string       `gorm:"unique;not null;;comment:邮箱"`
	Password string       `gorm:"type:varchar(60);not null;comment:密码"`
	Nickname *string      `gorm:"type:varchar(10);comment:昵称"`
	Gender   *int8        `gorm:"type:integer;comment:性别"`
	About    *string      `gorm:"type:varchar(60);comment:关于我"`
	Birthday sql.NullTime `gorm:"comment:生日"`
	Roles    []*Role      `gorm:"many2many:user_role;"`
}

func NewUserDAO() *UserDAO {
	return &UserDAO{BasicDAO: NewBasicDao()}
}

func (u *UserDAO) Insert(ctx context.Context, user *User) error {
	err := u.Datasource(ctx).Create(user).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return response.UserCreateDuplicateEmail
	}
	return err
}

func (u *UserDAO) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := u.Datasource(ctx).Preload("Roles").Where("email = ?", email).First(&user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return &user, response.UserLoginEmailNotFound
	}
	return &user, err
}

func (u *UserDAO) AssociateRoles(ctx context.Context, uid uint, roles *[]Role) error {
	return u.Datasource(ctx).Model(&User{
		Model: gorm.Model{
			ID: uid,
		},
	}).Association("Roles").Replace(roles)
}

func (u *UserDAO) FindByUid(ctx context.Context, uid uint, needRoles bool) (*User, error) {
	var user User
	query := u.Datasource(ctx).Model(&User{})
	if needRoles {
		query.Preload("Roles")
	}
	err := query.Where("id = ?", uid).First(&user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return &user, response.UserNotExist
	}
	return &user, err
}

func (u *UserDAO) FindRolesByUid(ctx context.Context, uid uint) ([]*Role, error) {
	var result []*Role
	err := u.Datasource(ctx).Table("sys_role as r").Select("r.id", "r.name").Joins("join sys_user_role as ur on r.id = ur.role_id and ur.user_id = ?", uid).Scan(&result).Error
	return result, err
}
