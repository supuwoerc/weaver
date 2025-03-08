package initialize

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"runtime/debug"
)

type CronLogger struct {
	logger *zap.SugaredLogger
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
					// FIXME:全局通用的邮件告警方法
					// TODO
					//go func() {
					//	if e := email.NewEmailClient().SendText(adminEmail, constant.CronRecover, message); e != nil {
					//		global.Logger.Errorf("发送邮件失败,信息:%s", e.Error())
					//	}
					//}()
				}
			}()
			j.Run()
		})
	}
}

// TODO:废弃
//func InitCron(logger *zap.SugaredLogger) (*cron.Cron, cron.Logger) {
//	l := &CronLogger{logger: logger}
//	return cron.New(cron.WithLogger(l), cron.WithSeconds(), cron.WithChain(CronRecover(l))), l
//}
