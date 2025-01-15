package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Email      string     `json:"email" orm:"unique;not null;;comment:邮箱"`
	Password   string     `json:"password,omitempty"`
	Nickname   *string    `json:"nickname"`
	Gender     *uint8     `json:"gender"`
	About      *string    `json:"about"`
	Birthday   *time.Time `json:"birthday"`
	Roles      []*Role    `json:"roles" gorm:"many2many:user_role;"`
	gorm.Model `json:"-"`
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
