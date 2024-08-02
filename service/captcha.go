package service

import (
	"gin-web/pkg/captcha"
	"github.com/gin-gonic/gin"
)

type CaptchaService struct {
	*BasicService
}

func NewCaptchaService(ctx *gin.Context) *CaptchaService {
	return &CaptchaService{
		BasicService: NewBasicService(ctx),
	}
}

func (c *CaptchaService) Generate() (*captcha.CaptchaInfo, error) {
	w := captcha.NewCaptcha()
	return w.Generate()
}

func (c *CaptchaService) Verify(id, answer string) bool {
	w := captcha.NewCaptcha()
	return w.Verify(id, answer)
}
