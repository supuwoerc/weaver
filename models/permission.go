package models

import "gin-web/pkg/database"

type Permission struct {
	Name      string  `json:"name" gorm:"not null;"`
	Resource  string  `json:"resource" gorm:"not null;"`
	Roles     []*Role `json:"roles" gorm:"many2many:role_permission;"`
	CreatorId uint    `json:"-" gorm:"not null;"`
	Creator   User    `json:"creator" gorm:"foreignKey:CreatorId;references:ID"`
	UpdaterId uint    `json:"-" gorm:"not null;"`
	Updater   User    `json:"updater" gorm:"foreignKey:UpdaterId;references:ID"`
	database.BasicModel
}
