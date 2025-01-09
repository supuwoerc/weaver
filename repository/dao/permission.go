package dao

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	Name     string  `gorm:"unique;not null;comment:权限名"`
	Resource string  `gorm:"unique;not null;comment:资源名"`
	Roles    []*Role `gorm:"many2many:role_permission;"`
}
