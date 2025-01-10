package captcha

import "github.com/mojocn/base64Captcha"

type CaptchaInfo struct {
	ID     string
	Base64 string
	Answer string
}

type Captcha struct {
	captchaClient *base64Captcha.Captcha
}

func NewCaptcha(height int, width int, length int, maxSkew float64, dotCount int) *Captcha {
	return &Captcha{captchaClient: base64Captcha.NewCaptcha(base64Captcha.NewDriverDigit(height, width, length, maxSkew, dotCount), &RedisStore{})}
}

func (c *Captcha) Generate() (*CaptchaInfo, error) {
	id, b64s, answer, err := c.captchaClient.Generate()
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
	return c.captchaClient.Verify(id, answer, true)
}
