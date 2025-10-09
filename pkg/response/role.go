package response

import (
	"github.com/samber/lo"
	"github.com/supuwoerc/weaver/models"
)

// RoleListRowResponse 角色列表的行
type RoleListRowResponse struct {
	*models.Role
	Users       any     `json:"users,omitempty"`       // 用户
	Permissions any     `json:"permissions,omitempty"` // 权限
	Creator     Creator `json:"creator"`               // 创建者
	Updater     Updater `json:"updater"`               // 更新人
}

// ToRoleListRowResponse 将role转为响应
func ToRoleListRowResponse(role *models.Role) *RoleListRowResponse {
	return &RoleListRowResponse{
		Role: role,
		Creator: Creator{
			User: &role.Creator,
		},
		Updater: Updater{
			User: &role.Updater,
		},
	}
}

// RoleDetailResponse 角色详情
type RoleDetailResponse struct {
	*models.Role
	Users       []*SimpleUser           `json:"users"`             // 用户
	Permissions []*RoleDetailPermission `json:"permissions"`       // 权限
	Creator     any                     `json:"creator,omitempty"` // 创建者
	Updater     any                     `json:"updater,omitempty"` // 更新者
}

type RoleDetailPermission struct {
	*models.Permission
	Roles   any `json:"roles,omitempty"`   // 角色
	Creator any `json:"creator,omitempty"` // 创建者
	Updater any `json:"updater,omitempty"` // 更新者
}

// ToRoleDetailResponse 将permission转为响应
func ToRoleDetailResponse(role *models.Role) *RoleDetailResponse {
	return &RoleDetailResponse{
		Role: role,
		Users: lo.Map(role.Users, func(item *models.User, _ int) *SimpleUser {
			return &SimpleUser{
				User: item,
			}
		}),
		Permissions: lo.Map(role.Permissions, func(item *models.Permission, _ int) *RoleDetailPermission {
			return &RoleDetailPermission{
				Permission: item,
			}
		}),
	}
}
