package response

type GetCaptchaResponse struct {
	ID     string `json:"id"`
	Base64 string `json:"base64"`
}
