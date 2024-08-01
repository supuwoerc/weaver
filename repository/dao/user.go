package dao

import (
	"context"
	"errors"
	"gin-web/pkg/constant"
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type UserDAO struct {
	*BasicDAO
}

var userDAO *UserDAO

type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null;;comment:邮箱"`
	Password string `gorm:"type:varchar(60);not null;comment:密码"`
	NickName string `gorm:"type:varchar(10);comment:昵称"`
	Gender   int    `gorm:"type:integer;comment:性别;default:0"`
	About    string `gorm:"type:varchar(60);comment:关于我"`
	Birthday int64  `gorm:"comment:生日"` // 生日
	Roles    []Role `gorm:"many2many:user_role"`
}

func NewUserDAO(ctx *gin.Context) *UserDAO {
	if userDAO == nil {
		userDAO = &UserDAO{BasicDAO: NewBasicDao(ctx)}
	}
	return userDAO
}

func (u *UserDAO) Insert(ctx context.Context, user User) error {
	err := u.db.WithContext(ctx).Create(&user).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return constant.GetError(u.ctx, response.USER_CREATE_DUPLICATE_EMAIL)
	}
	return err
}

func (u *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return user, constant.GetError(u.ctx, response.USER_LOGIN_EMAIL_NOT_FOUND)
	}
	return user, err
}

func (u *UserDAO) AssociateRoles(ctx context.Context, uid uint, roles []Role) error {
	return u.db.WithContext(ctx).Model(&User{
		Model: gorm.Model{
			ID: uid,
		},
	}).Association("Roles").Replace(roles)
}
