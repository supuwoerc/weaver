package response

type GetCaptchaResponse struct {
	ID     string `json:"id"`     // ID
	Base64 string `json:"base64"` // 图片base64
}
