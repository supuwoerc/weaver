package request

type CreateRoleRequest struct {
	Name string `json:"name" binding:"required"`
}

type GetRoleListRequest struct {
	Name   string `json:"name" form:"name" binding:"startswith=hi,contains=email"`
	Limit  int    `json:"limit" form:"limit" binding:"required,min=1,max=100"`
	Offset int    `json:"offset"  form:"limit" binding:"min=0"`
}
