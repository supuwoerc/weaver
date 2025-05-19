package models

import "github.com/supuwoerc/weaver/pkg/database"

type Role struct {
	Name        string        `json:"name" gorm:"not null"`
	Users       []*User       `json:"users" gorm:"many2many:user_role;"`
	Permissions []*Permission `json:"permissions" gorm:"many2many:role_permission;"`
	CreatorId   uint          `json:"-" gorm:"not null;"`
	Creator     User          `json:"creator" gorm:"foreignKey:CreatorId;references:ID"`
	UpdaterId   uint          `json:"-" gorm:"not null;"`
	Updater     User          `json:"updater" gorm:"foreignKey:UpdaterId;references:ID"`
	database.BasicModel
}
