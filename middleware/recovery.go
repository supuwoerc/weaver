package middleware

import (
	"gin-web/conf"
	"gin-web/pkg/constant"
	"gin-web/pkg/email"
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"runtime/debug"
)

type RecoveryMiddle struct {
	emailClient *email.EmailClient
	logger      *zap.SugaredLogger
	conf        *conf.Config
}

func NewRecoveryMiddleware(emailClient *email.EmailClient, logger *zap.SugaredLogger, conf *conf.Config) *RecoveryMiddle {
	return &RecoveryMiddle{
		emailClient: emailClient,
		logger:      logger,
		conf:        conf,
	}
}

func (r *RecoveryMiddle) Recovery() gin.HandlerFunc {
	adminEmail := r.conf.System.Admin.Email
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
