package models

import "gin-web/pkg/database"

type Permission struct {
	Name     string  `json:"name" gorm:"unique;not null;"`
	Resource string  `json:"resource" gorm:"unique;not null;"`
	Roles    []*Role `json:"roles" gorm:"many2many:role_permission;"`
	database.BasicModel
}
