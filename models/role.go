package models

import "gin-web/pkg/database"

type Role struct {
	Name        string        `json:"name" gorm:"not null"`
	Users       []*User       `json:"users" gorm:"many2many:user_role;"`
	Permissions []*Permission `json:"permissions" gorm:"many2many:role_permission;"`
	database.BasicModel
}
