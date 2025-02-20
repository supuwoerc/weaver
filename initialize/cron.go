package initialize

import (
	"fmt"
	"gin-web/pkg/email"
	"gin-web/pkg/global"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"runtime/debug"
)

type cronLogger struct {
	logger *zap.SugaredLogger
}

func (c *cronLogger) Info(msg string, keysAndValues ...interface{}) {
	c.logger.Infow(fmt.Sprintf("cron job:%s", msg), keysAndValues...)
}

func (c *cronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	c.logger.Errorw(fmt.Sprintf("cron job:%s", msg), append([]interface{}{"error", err}, keysAndValues...)...)
}

func cronRecover(logger cron.Logger) cron.JobWrapper {
	adminEmail := viper.GetString("system.admin.email")
	return func(j cron.Job) cron.Job {
		return cron.FuncJob(func() {
			defer func() {
				if r := recover(); r != nil {
					message := string(debug.Stack())
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					logger.Error(err, "panic", "stack", message)
					go func() {
						if e := email.SendText(adminEmail, "Recovery", message); e != nil {
							global.Logger.Errorf("发送邮件失败,信息:%s", e.Error())
						}
					}()
				}
			}()
			j.Run()
		})
	}
}

func InitCron(logger *zap.SugaredLogger) *cron.Cron {
	l := &cronLogger{logger: logger}
	return cron.New(cron.WithChain(cronRecover(l)), cron.WithLogger(l), cron.WithSeconds())
}
