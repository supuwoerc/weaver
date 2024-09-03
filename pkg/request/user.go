package request

type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	ID       string `json:"id" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SetRolesRequest struct {
	UserId  uint   `json:"user_id" binding:"required"`
	RoleIds []uint `json:"role_ids" binding:"required"`
}

type GetRolesRequest struct {
	UserId uint `form:"user_id" binding:"required"`
}
