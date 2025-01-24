package models

import (
	"gin-web/pkg/database"
	"time"
)

type User struct {
	Email    string     `json:"email" gorm:"unique;not null;;comment:邮箱"`
	Password string     `json:"password,omitempty"`
	Nickname *string    `json:"nickname"`
	Gender   *uint8     `json:"gender"`
	About    *string    `json:"about"`
	Birthday *time.Time `json:"birthday"`
	Roles    []*Role    `json:"roles" gorm:"many2many:user_role;"`
	database.BasicModel
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
