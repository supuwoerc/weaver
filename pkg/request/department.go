package request

// CreateDepartmentRequest 创建新部门的请求参数
type CreateDepartmentRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=20"`
	ParentId *uint  `json:"parent_id" binding:"omitempty,min=1"`
}
