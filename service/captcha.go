package service

import (
	"gin-web/pkg/captcha"
	"sync"
)

type CaptchaService struct {
	*BasicService
}

var (
	captchaOnce    sync.Once
	captchaService *CaptchaService
)

func NewCaptchaService() *CaptchaService {
	captchaOnce.Do(func() {
		captchaService = &CaptchaService{
			BasicService: NewBasicService(),
		}
	})
	return captchaService
}

func (c *CaptchaService) Generate() (*captcha.CaptchaInfo, error) {
	w := captcha.NewCaptcha()
	return w.Generate()
}

func (c *CaptchaService) Verify(id, answer string) bool {
	w := captcha.NewCaptcha()
	return w.Verify(id, answer)
}
