package request

type SignUpRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}
