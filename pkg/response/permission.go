package response

import (
	"github.com/supuwoerc/weaver/models"
	"github.com/supuwoerc/weaver/pkg/constant"
)

// PermissionListRowResponse 权限列表的行
type PermissionListRowResponse struct {
	*models.Permission
	Roles   any     `json:"roles,omitempty"` // 角色
	Creator Creator `json:"creator"`         // 创建者
	Updater Updater `json:"updater"`         // 更新者
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
	Creator any `json:"creator,omitempty"` // 创建者
	Updater any `json:"updater,omitempty"` // 更新者
}

type PermissionDetailRole struct {
	*models.Role
	Users       any `json:"users,omitempty"`       // 用户
	Permissions any `json:"permissions,omitempty"` // 权限
	Creator     any `json:"creator,omitempty"`     // 创建者
	Updater     any `json:"updater,omitempty"`     // 更新者
}

// ToPermissionDetailResponse 将permission转为响应
func ToPermissionDetailResponse(permission *models.Permission) *PermissionDetailResponse {
	return &PermissionDetailResponse{
		Permission: permission,
	}
}

type FrontEndPermissions []*FrontEndPermission

// FrontEndPermission 前端权限列表
type FrontEndPermission struct {
	ID       uint                    `json:"id"`       // ID
	Name     string                  `json:"name"`     // 权限名称
	Resource string                  `json:"resource"` // 资源
	Type     constant.PermissionType `json:"type"`     // 资源类型
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
