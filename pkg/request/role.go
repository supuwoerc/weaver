package request

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=20"`
	Users       []uint `json:"users" binding:"omitempty,dive,min=1"`
	Permissions []uint `json:"permissions" binding:"omitempty,dive,min=1"`
}

type GetRoleListRequest struct {
	Keyword string `json:"keyword" form:"keyword" binding:"omitempty,min=1,max=20"`
	Limit   int    `json:"limit" form:"limit" binding:"required,min=1,max=200"`
	Offset  int    `json:"offset"  form:"offset" binding:"min=0"`
}
