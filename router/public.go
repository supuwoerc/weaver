package router

import (
	v1 "gin-web/api/v1"
	"github.com/gin-gonic/gin"
)

func InitPublicRouter(r *gin.RouterGroup) {
	group := r.Group("public")
	{
		// 用户模块
		userApi := v1.NewUserApi()
		userGroup := group.Group("user")
		userGroup.POST("/signup", userApi.SignUp)
		userGroup.POST("/login", userApi.Login)
	}
	{
		// 验证码模块
		captchaApi := v1.NewCaptchaApi()
		captchaGroup := group.Group("captcha")
		captchaGroup.GET("/generate", captchaApi.GenerateCaptcha)
	}
}
