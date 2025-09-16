package models

import (
	"fmt"

	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/database"

	"github.com/samber/lo"
)

type Permission struct {
	Name      string                  `json:"name" gorm:"not null;"`
	Resource  string                  `json:"resource" gorm:"not null;"`
	Type      constant.PermissionType `json:"type" gorm:"not null;"`
	Roles     []*Role                 `json:"roles" gorm:"many2many:role_permission;"`
	CreatorID uint                    `json:"-" gorm:"not null;"`
	Creator   User                    `json:"creator" gorm:"foreignKey:CreatorID;references:ID"`
	UpdaterID uint                    `json:"-" gorm:"not null;"`
	Updater   User                    `json:"updater" gorm:"foreignKey:UpdaterID;references:ID"`
	database.BasicModel
}

func (p *Permission) GetRoleIds() []uint {
	return lo.Map(p.Roles, func(item *Role, index int) uint {
		return item.ID
	})
}

// GetResourceKey 获取资源的完整标识符
func (p *Permission) GetResourceKey() string {
	return fmt.Sprintf("%s:%s", p.Type.String(), p.Resource)
}

// IsApiPermission 判断是否为API权限
func (p *Permission) IsApiPermission() bool {
	return p.Type == constant.ApiRoute
}

// IsViewPermission 判断是否为视图权限
func (p *Permission) IsViewPermission() bool {
	return p.Type == constant.ViewRoute || p.Type == constant.ViewResource
}
