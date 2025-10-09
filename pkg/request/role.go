package request

// CreateRoleRequest 创建新角色的请求参数
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=20"`       // 角色名称
	Users       []uint `json:"users" binding:"omitempty,dive,min=1"`       // 角色关联用户
	Permissions []uint `json:"permissions" binding:"omitempty,dive,min=1"` // 角色关联权限
	ParentID    *uint  `json:"parent_id" binding:"omitempty,min=1"`        // 父角色ID
}

// GetRoleListRequest 查询角色列表的参数
type GetRoleListRequest struct {
	Keyword string `json:"keyword" form:"keyword" binding:"omitempty,min=1,max=20"` // 关键字
	Limit   int    `json:"limit" form:"limit" binding:"required,min=1,max=200"`     // 分页数量
	Offset  int    `json:"offset"  form:"offset" binding:"min=0"`                   // 分页偏移
}

// GetRoleDetailRequest 查询角色详情的参数
type GetRoleDetailRequest struct {
	ID uint `json:"id" form:"id" binding:"required,min=1"` // ID
}

// DeleteRoleRequest 删除角色的参数
type DeleteRoleRequest = GetRoleDetailRequest

// UpdateRoleRequest 更新角色的参数
type UpdateRoleRequest struct {
	ID uint `json:"id" form:"id" binding:"required,min=1"` // ID
	CreateRoleRequest
}
