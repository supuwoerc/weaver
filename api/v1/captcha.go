package v1

import (
	"gin-web/pkg/response"
	"gin-web/service"
	"github.com/gin-gonic/gin"
	"sync"
)

type CaptchaApi struct {
	*BasicApi
	service *service.CaptchaService
}

var (
	captchaOnce sync.Once
	captchaApi  *CaptchaApi
)

func NewCaptchaApi() *CaptchaApi {
	captchaOnce.Do(func() {
		captchaApi = &CaptchaApi{
			BasicApi: NewBasicApi(),
			service:  service.NewCaptchaService(),
		}
	})
	return captchaApi
}

func (c *CaptchaApi) GenerateCaptcha(ctx *gin.Context) {
	captchaInfo, err := c.service.Generate(service.SignUp)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.SuccessWithData[response.GetCaptchaResponse](ctx, response.GetCaptchaResponse{
		ID:     captchaInfo.ID,
		Base64: captchaInfo.Base64,
	})
}
