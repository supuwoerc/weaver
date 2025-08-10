package captcha

import (
	v1 "github.com/supuwoerc/weaver/api/v1"
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/response"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Generate(t constant.CaptchaType) (*response.GetCaptchaResponse, error)
}

type Api struct {
	*v1.BasicApi
	service Service
}

func NewCaptchaApi(basic *v1.BasicApi, service Service) *Api {
	captchaApi := &Api{
		BasicApi: basic,
		service:  service,
	}
	captchaGroup := basic.Route.Group("public/captcha")
	{
		captchaGroup.GET("signup", captchaApi.GenerateSignUpCaptcha)
	}
	return captchaApi
}

func (c *Api) commonGenerate(ctx *gin.Context, t constant.CaptchaType) {
	res, err := c.service.Generate(t)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.SuccessWithData(ctx, res)
}

// GenerateSignUpCaptcha 注册验证码
//
//	@Summary		生成注册验证码
//	@Description	生成用户注册时使用的验证码
//	@Tags			验证码管理
//	@Accept			json
//	@Produce		json
//	@Success		10000	{object}	response.BasicResponse[response.GetCaptchaResponse]	"生成成功，code=10000"
//	@Failure		10001	{object}	response.BasicResponse[any]							"生成失败，code=10001"
//	@Router			/public/captcha/signup [get]
func (c *Api) GenerateSignUpCaptcha(ctx *gin.Context) {
	c.commonGenerate(ctx, constant.SignUp)
}
