package request

// CreateDepartmentRequest 创建新部门的请求参数
type CreateDepartmentRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=20"`   // 部门名称
	ParentID *uint  `json:"parent_id" binding:"omitempty,min=1"`    // 父部门ID
	Leaders  []uint `json:"leaders" binding:"omitempty,dive,min=1"` // 部门leader集合
	Users    []uint `json:"users" binding:"omitempty,dive,min=1"`   // 部门用户ID集合
}

// GetDepartmentTreeRequest 查询组织架构树的参数
type GetDepartmentTreeRequest struct {
	WithCrew bool `form:"with_crew" binding:"boolean"` // 是否返回人员信息
}

// GetDepartmentsByParentIDRequest 查询组织架构树的参数
type GetDepartmentsByParentIDRequest struct {
	ParentID *uint `json:"parent_id" binding:"omitempty,min=1"` // 父部门ID
	WithCrew bool  `form:"with_crew" binding:"boolean"`         // 是否返回人员信息
}
