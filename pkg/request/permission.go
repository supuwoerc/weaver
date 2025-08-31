package request

import "github.com/supuwoerc/weaver/pkg/constant"

// CreatePermissionRequest 创建新权限的请求参数
type CreatePermissionRequest struct {
	Name     string                  `json:"name" binding:"required,min=1,max=20"`
	Resource string                  `json:"resource" binding:"required,min=1,max=255"`
	Type     constant.PermissionType `json:"type" binding:"required,oneof=1 2 3 4"`
	Roles    []uint                  `json:"roles" binding:"omitempty,dive,min=1"`
}

// GetPermissionListRequest 查询权限列表的参数
type GetPermissionListRequest struct {
	Keyword string `json:"keyword" form:"keyword" binding:"omitempty,min=1,max=20"`
	Limit   int    `json:"limit" form:"limit" binding:"required,min=1,max=200"`
	Offset  int    `json:"offset"  form:"offset" binding:"min=0"`
}

// GetPermissionDetailRequest 查询权限详情的参数
type GetPermissionDetailRequest struct {
	ID uint `json:"id" form:"id" binding:"required,min=1"`
}

// DeletePermissionRequest 删除权限的参数
type DeletePermissionRequest = GetPermissionDetailRequest

// UpdatePermissionRequest 更新权限的参数
type UpdatePermissionRequest struct {
	ID uint `json:"id" form:"id" binding:"required,min=1"`
	CreatePermissionRequest
}
