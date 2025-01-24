package models

import (
	"gin-web/pkg/database"
	"time"
)

type User struct {
	Email    string      `json:"email" gorm:"unique;not null;;comment:邮箱"`
	Password string      `json:"-"`
	Nickname *string     `json:"nickname"`
	AvatarId *uint       `json:"-"`
	Avatar   *Attachment `json:"avatar" gorm:"foreignKey:AvatarId;references:ID"`
	Gender   *uint8      `json:"gender"`
	About    *string     `json:"about"`
	Birthday *time.Time  `json:"birthday"`
	Roles    []*Role     `json:"roles" gorm:"many2many:user_role;"`
	database.BasicModel
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
