package models

import "gin-web/pkg/database"

type Department struct {
	Name      string      `json:"name" gorm:"not null;"`
	ParentId  *uint       `json:"-"` // 父级部门ID,nil则为顶级部门
	Parent    *Department `json:"parent" gorm:"foreignKey:ParentId;references:ID"`
	Ancestors *string     `json:"-"` // 祖先部门路径逗号拼接的字符串
	Leaders   []*User     `json:"leaders" gorm:"many2many:department_leader;"`
	Users     []*User     `json:"users" gorm:"many2many:user_department;"`
	CreatorId uint        `json:"-" gorm:"not null;"`
	Creator   User        `json:"creator" gorm:"foreignKey:CreatorId;references:ID"`
	UpdaterId uint        `json:"-" gorm:"not null;"`
	Updater   User        `json:"updater" gorm:"foreignKey:UpdaterId;references:ID"`
	database.BasicModel
}
