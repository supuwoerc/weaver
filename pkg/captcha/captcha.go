package captcha

import "github.com/mojocn/base64Captcha"

var captchaClient *base64Captcha.Captcha

type CaptchaInfo struct {
	ID     string
	Base64 string
	Answer string
}

type Captcha struct {
}

func init() {
	captchaClient = base64Captcha.NewCaptcha(base64Captcha.NewDriverDigit(100, 380, 6, 0.3, 80), RedisStore{})
}

var captcha *Captcha

func NewCaptcha() *Captcha {
	if captcha == nil {
		captcha = &Captcha{}
	}
	return captcha
}

func (c *Captcha) Generate() (*CaptchaInfo, error) {
	id, b64s, answer, err := captchaClient.Generate()
	if err != nil {
		return nil, err
	}
	return &CaptchaInfo{
		ID:     id,
		Base64: b64s,
		Answer: answer,
	}, nil
}

func (c *Captcha) Verify(id, answer string) bool {
	return captchaClient.Verify(id, answer, true)
}
