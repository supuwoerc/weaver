package models

import "gorm.io/gorm"

type Permission struct {
	Name       string  `json:"name" gorm:"unique;not null;"`
	Resource   string  `json:"resource" gorm:"unique;not null;"`
	Roles      []*Role `json:"roles,omitempty" gorm:"many2many:role_permission;"`
	gorm.Model `json:"-"`
}
