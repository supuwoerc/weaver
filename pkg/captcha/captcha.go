package captcha

import "github.com/mojocn/base64Captcha"

type CommonCaptchaInfo struct {
	ID     string
	Base64 string
	Answer string
}

type Captcha struct {
	captchaClient *base64Captcha.Captcha
}

func NewCaptcha(height int, width int, length int, maxSkew float64, dotCount int, store base64Captcha.Store) *Captcha {
	return &Captcha{
		captchaClient: base64Captcha.NewCaptcha(base64Captcha.NewDriverDigit(height, width, length, maxSkew, dotCount),
			store,
		),
	}
}

func (c *Captcha) Generate() (*CommonCaptchaInfo, error) {
	id, b64s, answer, err := c.captchaClient.Generate()
	if err != nil {
		return nil, err
	}
	return &CommonCaptchaInfo{
		ID:     id,
		Base64: b64s,
		Answer: answer,
	}, nil
}

func (c *Captcha) Verify(id, answer string) bool {
	return c.captchaClient.Verify(id, answer, true)
}
