package models

import "gorm.io/gorm"

type Role struct {
	Name        string        `json:"name" gorm:"unique;not null"`
	Users       []*User       `json:"users,omitempty" gorm:"many2many:user_role;"`
	Permissions []*Permission `json:"permissions,omitempty" gorm:"many2many:role_permission;"`
	gorm.Model  `json:"-"`
}
