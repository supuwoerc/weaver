package constant

//go:generate stringer -type=CaptchaType -linecomment -output captcha_string.go
type CaptchaType int

const (
	Default CaptchaType = iota + 1 // defaultCaptcha
	SignUp                         // signupCaptcha
)
