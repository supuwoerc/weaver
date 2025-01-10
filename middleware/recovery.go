package middleware

import (
	"gin-web/pkg/email"
	"gin-web/pkg/global"
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"runtime/debug"
)

func Recovery() gin.HandlerFunc {
	adminEmail := viper.GetString("system.admin.email")
	return gin.CustomRecovery(func(c *gin.Context, err any) {
		message := string(debug.Stack())
		global.Logger.Errorf("Recovery panic,堆栈信息:%s", message)
		go func() {
			if e := email.SendText(adminEmail, "Recovery", message); e != nil {
				global.Logger.Errorf("发送邮件失败,信息:%s", e.Error())
			}
		}()
		response.HttpResponse[any](c, response.RecoveryError, nil, nil, nil)
	})
}
