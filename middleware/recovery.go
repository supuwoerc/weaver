package middleware

import (
	"gin-web/pkg/constant"
	"gin-web/pkg/email"
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"runtime/debug"
	"sync"
)

type RecoveryMiddle struct {
	emailClient *email.EmailClient
	logger      *zap.SugaredLogger
	viper       *viper.Viper
}

var (
	recoveryOnce   sync.Once
	recoveryMiddle *RecoveryMiddle
)

// TODO:确认是否需要单例
func NewRecoveryMiddleware(emailClient *email.EmailClient, logger *zap.SugaredLogger, v *viper.Viper) *RecoveryMiddle {
	recoveryOnce.Do(func() {
		recoveryMiddle = &RecoveryMiddle{
			emailClient: emailClient,
			logger:      logger,
			viper:       v,
		}
	})
	return recoveryMiddle
}

func (r *RecoveryMiddle) Recovery() gin.HandlerFunc {
	adminEmail := r.viper.GetString("system.admin.email")
	return gin.CustomRecovery(func(c *gin.Context, err any) {
		message := string(debug.Stack())
		r.logger.Errorf("Recovery panic,堆栈信息:%s", message)
		// TODO:全局通用的告警方法
		go func() {
			if e := r.emailClient.SendText(adminEmail, constant.Recover, message); e != nil {
				r.logger.Errorf("发送邮件失败,信息:%s", e.Error())
			}
		}()
		response.HttpResponse[any](c, response.RecoveryError, nil, nil, nil)
	})
}
