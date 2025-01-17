package request

type CreatePermissionRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=20"`
	Resource string `json:"resource" binding:"required,min=1,max=255"`
	Roles    []uint `json:"roles" binding:"omitempty,dive,min=1"`
}

type GetPermissionListRequest struct {
	Keyword string `json:"keyword" form:"keyword" binding:"omitempty,min=1,max=20"`
	Limit   int    `json:"limit" form:"limit" binding:"required,min=1,max=200"`
	Offset  int    `json:"offset"  form:"offset" binding:"min=0"`
}
