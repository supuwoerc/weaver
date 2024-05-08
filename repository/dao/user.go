package dao

import (
	"context"
	"errors"
	"gin-web/pkg/constant"
	"gin-web/pkg/global"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type UserDAO struct {
	db *gorm.DB
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
}

func NewUserDAO() *UserDAO {
	if userDAO == nil {
		userDAO = &UserDAO{db: global.DB}
	}
	return userDAO
}

func (u UserDAO) Insert(ctx context.Context, user User) error {
	err := u.db.WithContext(ctx).Create(&user).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return constant.USER_CREATE_DUPLICATE_EMAIL_ERR
	}
	return err
}
