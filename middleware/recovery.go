package middleware

import (
	"runtime/debug"

	"github.com/supuwoerc/weaver/conf"
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/logger"
	"github.com/supuwoerc/weaver/pkg/response"

	"github.com/gin-gonic/gin"
)

type RecoverEmailClient interface {
	Alarm2Admin(subject constant.Subject, body string) error
}

type RecoveryMiddle struct {
	emailClient RecoverEmailClient
	logger      *logger.Logger
	conf        *conf.Config
}

func NewRecoveryMiddleware(emailClient RecoverEmailClient, logger *logger.Logger, conf *conf.Config) *RecoveryMiddle {
	return &RecoveryMiddle{
		emailClient: emailClient,
		logger:      logger,
		conf:        conf,
	}
}

func (r *RecoveryMiddle) Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err any) {
		message := string(debug.Stack())
		r.logger.Errorf("recover: %v,recovery panic,stack info: %s", err, message)
		go func() {
			if e := r.emailClient.Alarm2Admin(constant.Recover, message); e != nil {
				r.logger.Errorf("send revover email to admin fail,stack info: %s", e.Error())
			}
		}()
		response.HttpResponse[any](c, response.RecoveryError, nil, nil, nil)
	})
}
