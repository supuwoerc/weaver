package request

// CreateDepartmentRequest 创建新部门的请求参数
type CreateDepartmentRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=20"`
	ParentId *uint  `json:"parent_id" binding:"omitempty,min=1"`
	Leaders  []uint `json:"leaders" binding:"omitempty,dive,min=1"`
	Users    []uint `json:"users" binding:"omitempty,dive,min=1"`
}

// GetDepartmentTreeRequest 查询组织架构树的参数
type GetDepartmentTreeRequest struct {
	WithCrew bool `form:"with_crew" binding:"boolean"`
}
