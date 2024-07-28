package service

import (
	"gin-web/pkg/captcha"
	"github.com/gin-gonic/gin"
)

type CaptchaService struct {
	*BasicService
}

var captchaService *CaptchaService

func NewCaptchaService(ctx *gin.Context) *CaptchaService {
	if captchaService == nil {
		captchaService = &CaptchaService{
			BasicService: NewBasicService(ctx),
		}
	}
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
