package response

import (
	"github.com/samber/lo"
	"github.com/supuwoerc/weaver/models"
)

// RoleListRowResponse 角色列表的行
type RoleListRowResponse struct {
	*models.Role
	Users       any     `json:"users,omitempty"`
	Permissions any     `json:"permissions,omitempty"`
	Creator     Creator `json:"creator"`
	Updater     Updater `json:"updater"`
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
	Users       []*SimpleUser           `json:"users"`
	Permissions []*RoleDetailPermission `json:"permissions"`
	Creator     any                     `json:"creator,omitempty"`
	Updater     any                     `json:"updater,omitempty"`
}

type RoleDetailPermission struct {
	*models.Permission
	Roles   any `json:"roles,omitempty"`
	Creator any `json:"creator,omitempty"`
	Updater any `json:"updater,omitempty"`
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
