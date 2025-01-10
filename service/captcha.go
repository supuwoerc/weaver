package service

import (
	"gin-web/pkg/captcha"
	"sync"
)

type CaptchaService struct {
	*BasicService
	client *captcha.Captcha
}

var (
	captchaOnce    sync.Once
	captchaService *CaptchaService
)

func NewCaptchaService() *CaptchaService {
	captchaOnce.Do(func() {
		captchaService = &CaptchaService{
			BasicService: NewBasicService(),
			client:       captcha.NewCaptcha(100, 348, 6, 0.3, 80),
		}
	})
	return captchaService
}

func (c *CaptchaService) Generate() (*captcha.CaptchaInfo, error) {
	return c.client.Generate()
}

func (c *CaptchaService) Verify(id, answer string) bool {
	return c.client.Verify(id, answer)
}
