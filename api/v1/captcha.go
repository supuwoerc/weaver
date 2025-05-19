package v1

import (
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/response"

	"github.com/gin-gonic/gin"
)

type CaptchaService interface {
	Generate(t constant.CaptchaType) (*response.GetCaptchaResponse, error)
}

type CaptchaApi struct {
	service CaptchaService
}

func NewCaptchaApi(route *gin.RouterGroup, service CaptchaService) *CaptchaApi {
	captchaApi := &CaptchaApi{
		service: service,
	}
	// 挂载路由
	captchaGroup := route.Group("public/captcha")
	{
		captchaGroup.GET("signup", captchaApi.GenerateSignUpCaptcha)
	}
	return captchaApi
}

func (c *CaptchaApi) commonGenerate(ctx *gin.Context, t constant.CaptchaType) {
	res, err := c.service.Generate(t)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.SuccessWithData(ctx, res)
}

// GenerateSignUpCaptcha 注册验证码
func (c *CaptchaApi) GenerateSignUpCaptcha(ctx *gin.Context) {
	c.commonGenerate(ctx, constant.SignUp)
}
