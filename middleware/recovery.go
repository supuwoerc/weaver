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
	emailClient *email.Client
	logger      *zap.SugaredLogger
	conf        *conf.Config
}

func NewRecoveryMiddleware(emailClient *email.Client, logger *zap.SugaredLogger, conf *conf.Config) *RecoveryMiddle {
	return &RecoveryMiddle{
		emailClient: emailClient,
		logger:      logger,
		conf:        conf,
	}
}

func (r *RecoveryMiddle) Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err any) {
		message := string(debug.Stack())
		r.logger.Errorf("Recovery panic,堆栈信息:%s", message)
		go func() {
			if e := r.emailClient.Alarm2Admin(constant.Recover, message); e != nil {
				r.logger.Errorf("发送邮件失败,信息:%s", e.Error())
			}
		}()
		response.HttpResponse[any](c, response.RecoveryError, nil, nil, nil)
	})
}
