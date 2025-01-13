package request

type CreateRoleRequest struct {
	Name string `json:"name" binding:"required"`
}

type GetRoleListRequest struct {
	Name   string `json:"name"`
	Limit  int    `json:"limit" binding:"required"`
	Offset int    `json:"offset"`
}
