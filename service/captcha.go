package service

import (
	"gin-web/pkg/captcha"
	"gin-web/pkg/response"
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

func (c *CaptchaService) Generate(t CaptchaType) (*response.GetCaptchaResponse, error) {
	var info *captcha.CaptchaInfo
	var err error
	if target, ok := c.clients[t]; ok {
		info, err = target.Generate()
	} else {
		info, err = c.clients[Default].Generate()
	}
	if err != nil {
		return nil, err
	}
	return &response.GetCaptchaResponse{
		ID:     info.ID,
		Base64: info.Base64,
	}, nil
}

func (c *CaptchaService) Verify(t CaptchaType, id, answer string) bool {
	if target, ok := c.clients[t]; ok {
		return target.Verify(id, answer)
	} else {
		return c.clients[Default].Verify(id, answer)
	}
}
