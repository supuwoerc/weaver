package service

import (
	"gin-web/pkg/captcha"
	"gin-web/pkg/constant"
	"gin-web/pkg/response"
)

type CaptchaService struct {
	clients map[constant.CaptchaType]*captcha.Captcha // 不同业务使用不同参数的验证码生成器
}

func NewCaptchaService() *CaptchaService {
	return &CaptchaService{
		clients: map[constant.CaptchaType]*captcha.Captcha{
			constant.Default: captcha.NewCaptcha(100, 200, 6, 0.3, 80), // 默认验证码
			constant.SignUp:  captcha.NewCaptcha(100, 348, 6, 0.3, 80), // 注册验证码
		},
	}
}

func (c *CaptchaService) Generate(t constant.CaptchaType) (*response.GetCaptchaResponse, error) {
	var info *captcha.CommonCaptchaInfo
	var err error
	if target, ok := c.clients[t]; ok {
		info, err = target.Generate()
	} else {
		info, err = c.clients[constant.Default].Generate()
	}
	if err != nil {
		return nil, err
	}
	return &response.GetCaptchaResponse{
		ID:     info.ID,
		Base64: info.Base64,
	}, nil
}

func (c *CaptchaService) Verify(t constant.CaptchaType, id, answer string) bool {
	if target, ok := c.clients[t]; ok {
		return target.Verify(id, answer)
	} else {
		return c.clients[constant.Default].Verify(id, answer)
	}
}
