package v1

import (
	"gin-web/pkg/constant"
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
	"sync"
)

type CaptchaService interface {
	Generate(t constant.CaptchaType) (*response.GetCaptchaResponse, error)
}

type CaptchaApi struct {
	service CaptchaService
}

var (
	captchaOnce sync.Once
	captchaApi  *CaptchaApi
)

func NewCaptchaApi(route *gin.RouterGroup, service CaptchaService) *CaptchaApi {
	captchaOnce.Do(func() {
		captchaApi = &CaptchaApi{
			service: service,
		}
		// 挂载路由
		captchaGroup := route.Group("captcha")
		{
			captchaGroup.GET("generate", captchaApi.GenerateCaptcha)
		}
	})
	return captchaApi
}

func (c *CaptchaApi) GenerateCaptcha(ctx *gin.Context) {
	res, err := c.service.Generate(constant.SignUp)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.SuccessWithData(ctx, res)
}
