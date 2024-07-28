package v1

import (
	"gin-web/pkg/response"
	"gin-web/service"
	"github.com/gin-gonic/gin"
)

type CaptchaApi struct {
	*BasicApi
	service func(ctx *gin.Context) *service.CaptchaService
}

func NewCaptchaApi() CaptchaApi {
	return CaptchaApi{
		BasicApi: NewBasicApi(),
		service: func(ctx *gin.Context) *service.CaptchaService {
			return service.NewCaptchaService(ctx)
		},
	}
}

func (c CaptchaApi) GenerateCaptcha(ctx *gin.Context) {
	captchaInfo, err := c.service(ctx).Generate()
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.SuccessWithData[response.GetCaptchaResponse](ctx, response.GetCaptchaResponse{
		ID:     captchaInfo.ID,
		Base64: captchaInfo.Base64,
	})
}
