package request

// SignUpRequest 注册请求参数
type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email,max=50"`
	Password string `json:"password" binding:"required"`
	ID       string `json:"id" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=50"`
	Password string `json:"password" binding:"required"`
}

// GetUserListRequest 查询用户列表的参数
type GetUserListRequest struct {
	Keyword string `json:"keyword" form:"keyword" binding:"omitempty,min=1,max=20"`
	Limit   int    `json:"limit" form:"limit" binding:"required,min=1,max=200"`
	Offset  int    `json:"offset"  form:"offset" binding:"min=0"`
}

// ActiveAccountRequest 激活账户的请求参数
type ActiveAccountRequest struct {
	ActiveCode string `json:"active_code" form:"active_code" binding:"required,len=16"`
	ID         uint   `json:"id" form:"id" binding:"required,min=1"`
}
