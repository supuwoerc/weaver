package request

type CreateRoleRequest struct {
	Name string `json:"name" binding:"required"`
}

type GetRoleListRequest struct {
	Name   string `json:"name" form:"name"`
	Limit  int    `json:"limit" form:"limit" binding:"required,min=1,max=200"`
	Offset int    `json:"offset"  form:"offset" binding:"min=0"`
}
