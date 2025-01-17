package service

import (
	"gin-web/pkg/captcha"
	"sync"
)

//go:generate stringer -type=CaptchaType -linecomment -output captcha_string.go
type CaptchaType int

const (
	Default CaptchaType = iota + 1 // 默认验证码
	SignUp                         // 注册验证码
)

type CaptchaService struct {
	*BasicService
	clients map[CaptchaType]*captcha.Captcha // 不同业务使用不同参数的验证码生成器
}

var (
	captchaOnce    sync.Once
	captchaService *CaptchaService
)

func NewCaptchaService() *CaptchaService {
	captchaOnce.Do(func() {
		captchaService = &CaptchaService{
			BasicService: NewBasicService(),
			clients: map[CaptchaType]*captcha.Captcha{
				Default: captcha.NewCaptcha(100, 200, 6, 0.3, 80), // 默认验证码
				SignUp:  captcha.NewCaptcha(100, 348, 6, 0.3, 80), // 注册验证码
			},
		}
	})
	return captchaService
}

func (c *CaptchaService) Generate(t CaptchaType) (*captcha.CaptchaInfo, error) {
	if target, ok := c.clients[t]; ok {
		return target.Generate()
	} else {
		return c.clients[Default].Generate()
	}
}

func (c *CaptchaService) Verify(t CaptchaType, id, answer string) bool {
	if target, ok := c.clients[t]; ok {
		return target.Verify(id, answer)
	} else {
		return c.clients[Default].Verify(id, answer)
	}
}
