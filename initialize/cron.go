package initialize

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/logger"

	"github.com/robfig/cron/v3"
)

type CronLogger struct {
	logger      *logger.Logger
	emailClient *EmailClient
}

func (c *CronLogger) Info(msg string, keysAndValues ...interface{}) {
	c.logger.Infow(fmt.Sprintf("cron job: %s", msg), keysAndValues...)
}

func (c *CronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	c.logger.Errorw(fmt.Sprintf("cron job: %s", msg), append([]interface{}{"error", err}, keysAndValues...)...)
}
func (c *CronLogger) CronRecover() cron.JobWrapper {
	return func(j cron.Job) cron.Job {
		return cron.FuncJob(func() {
			defer func() {
				if r := recover(); r != nil {
					message := string(debug.Stack())
					c.logger.Errorw("cron recover", "panic", r, "stack", message)
					go func() {
						if e := c.emailClient.Alarm2Admin(context.Background(), constant.CronRecover, message); e != nil {
							c.logger.Errorw("cron recover alarm to admin fail", "err", e.Error())
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

func NewCronLogger(logger *logger.Logger, emailClient *EmailClient) *CronLogger {
	return &CronLogger{logger: logger, emailClient: emailClient}
}
