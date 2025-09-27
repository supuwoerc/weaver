package models

import "github.com/supuwoerc/weaver/pkg/database"

type Role struct {
	Name        string        `json:"name" gorm:"not null"`
	Users       []*User       `json:"users" gorm:"many2many:user_role;"`
	Permissions []*Permission `json:"permissions" gorm:"many2many:role_permission;"`
	ParentID    *uint         `json:"parent_id"`
	Parent      *Role         `json:"parent" gorm:"foreignKey:ParentID;references:ID"`
	Children    []*Role       `json:"children" gorm:"foreignKey:ParentID"`
	Ancestors   *string       `json:"ancestors"`
	CreatorID   uint          `json:"-" gorm:"not null;"`
	Creator     User          `json:"creator" gorm:"foreignKey:CreatorID;references:ID"`
	UpdaterID   uint          `json:"-" gorm:"not null;"`
	Updater     User          `json:"updater" gorm:"foreignKey:UpdaterID;references:ID"`
	database.BasicModel
}
