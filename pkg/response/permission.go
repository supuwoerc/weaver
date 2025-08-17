package response

import (
	"github.com/samber/lo"
	"github.com/supuwoerc/weaver/models"
	"github.com/supuwoerc/weaver/pkg/constant"
)

// PermissionListRowResponse 权限列表的行
type PermissionListRowResponse struct {
	*models.Permission
	Roles   any     `json:"roles,omitempty"`
	Creator Creator `json:"creator"`
	Updater Updater `json:"updater"`
}

// ToPermissionListRowResponse 将permission转为响应
func ToPermissionListRowResponse(permission *models.Permission) *PermissionListRowResponse {
	return &PermissionListRowResponse{
		Permission: permission,
		Creator: Creator{
			User: &permission.Creator,
		},
		Updater: Updater{
			User: &permission.Updater,
		},
	}
}

// PermissionDetailResponse 权限详情
type PermissionDetailResponse struct {
	*models.Permission
	Roles   []*PermissionDetailRole `json:"roles"`
	Creator any                     `json:"creator,omitempty"`
	Updater any                     `json:"updater,omitempty"`
}

type PermissionDetailRole struct {
	*models.Role
	Users       any `json:"users,omitempty"`
	Permissions any `json:"permissions,omitempty"`
	Creator     any `json:"creator,omitempty"`
	Updater     any `json:"updater,omitempty"`
}

// ToPermissionDetailResponse 将permission转为响应
func ToPermissionDetailResponse(permission *models.Permission) *PermissionDetailResponse {
	return &PermissionDetailResponse{
		Permission: permission,
		Roles: lo.Map(permission.Roles, func(item *models.Role, _ int) *PermissionDetailRole {
			return &PermissionDetailRole{
				Role: item,
			}
		}),
	}
}

type FrontEndPermissions []*FrontEndPermission

// FrontEndPermission 前端权限列表
type FrontEndPermission struct {
	ID       uint                    `json:"id"`
	Name     string                  `json:"name"`
	Resource string                  `json:"resource"`
	Type     constant.PermissionType `json:"type"`
}

// ToFrontEndPermissionResponse 将permission转为响应
func ToFrontEndPermissionResponse(permission *models.Permission) *FrontEndPermission {
	return &FrontEndPermission{
		ID:       permission.ID,
		Name:     permission.Name,
		Resource: permission.Resource,
		Type:     permission.Type,
	}
}
