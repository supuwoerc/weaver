package dao

import (
	"context"
	"gin-web/pkg/global"
	"gorm.io/gorm"
)

type UserDAO struct {
	db *gorm.DB
}

var userDAO *UserDAO

type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null;;comment:邮箱"`
	Password string `gorm:"type:varchar(50);not null;comment:密码"`
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
	return u.db.WithContext(ctx).Create(&user).Error
}
