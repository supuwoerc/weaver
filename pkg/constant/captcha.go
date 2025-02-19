package constant

//go:generate stringer -type=CaptchaType -linecomment -output captcha_string.go
type CaptchaType int

const (
	Default CaptchaType = iota + 1 // 默认验证码
	SignUp                         // 注册验证码
)
