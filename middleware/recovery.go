package middleware

import (
	"context"
	"runtime/debug"

	"github.com/supuwoerc/weaver/conf"
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/logger"
	"github.com/supuwoerc/weaver/pkg/response"

	"github.com/gin-gonic/gin"
)

type RecoverEmailClient interface {
	Alarm2Admin(ctx context.Context, subject constant.Subject, body string) error
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
		r.logger.WithContext(c).Errorw("recover panic", "panic", err, "stack", message)
		// copy context with context inner Keys
		copyCtx := c.Copy()
		go func(ctx context.Context) {
			if e := r.emailClient.Alarm2Admin(ctx, constant.Recover, message); e != nil {
				r.logger.WithContext(ctx).Errorf("send revover email to admin fail,stack info: %s", e.Error())
			}
		}(copyCtx)
		response.HttpResponse[any](c, response.RecoveryError, nil, nil, nil)
	})
}
