package models

import (
	"gin-web/pkg/constant"
	"gin-web/pkg/database"
	"time"
)

type User struct {
	Email       string               `json:"email" gorm:"unique;not null;"`
	Password    string               `json:"-"`
	Status      constant.UserStatus  `json:"status"`
	Nickname    *string              `json:"nickname"`
	Avatar      *string              `json:"avatar"`
	Gender      *constant.UserGender `json:"gender"`
	About       *string              `json:"about"`
	Birthday    *time.Time           `json:"birthday"`
	Roles       []*Role              `json:"roles" gorm:"many2many:user_role;"`
	Departments []*Department        `json:"departments" gorm:"many2many:user_department;"`
	database.BasicModel
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
