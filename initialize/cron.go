package initialize

import (
	"fmt"
	"runtime/debug"

	"github.com/supuwoerc/weaver/pkg/constant"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type CronLogger struct {
	logger      *zap.SugaredLogger
	emailClient *EmailClient
}

func (c *CronLogger) Info(msg string, keysAndValues ...interface{}) {
	c.logger.Infow(fmt.Sprintf("cron job:%s", msg), keysAndValues...)
}

func (c *CronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	c.logger.Errorw(fmt.Sprintf("cron job:%s", msg), append([]interface{}{"error", err}, keysAndValues...)...)
}
func (c *CronLogger) CronRecover() cron.JobWrapper {
	//adminEmail := viper.GetString("system.admin.email")
	return func(j cron.Job) cron.Job {
		return cron.FuncJob(func() {
			defer func() {
				if r := recover(); r != nil {
					message := string(debug.Stack())
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					c.logger.Error(err, "panic", "stack", message)
					go func() {
						if e := c.emailClient.Alarm2Admin(constant.CronRecover, message); e != nil {
							c.logger.Errorf("cron recover alarm fail:%s", e.Error())
						}
					}()
				}
			}()
			j.Run()
		})
	}
}

func NewCronClient(l *CronLogger) *cron.Cron {
	return cron.New(cron.WithLogger(l), cron.WithSeconds(), cron.WithChain(l.CronRecover()))
}

func NewCronLogger(logger *zap.SugaredLogger, emailClient *EmailClient) *CronLogger {
	return &CronLogger{logger: logger, emailClient: emailClient}
}
