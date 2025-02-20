package bootstrap

import (
	"gin-web/pkg/constant"
	"gin-web/pkg/job"
	"github.com/robfig/cron/v3"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	jobs    []job.SystemJob
	mapping = make(map[string]cron.EntryID)
)

func init() {
	jobs = []job.SystemJob{
		job.NewServerStatus(),
	}
}

func RegisterJobs(c *cron.Cron, logger *zap.SugaredLogger) error {
	key := "system.hooks.launch"
	onLaunch := viper.GetStringSlice(key)
	if lo.Contains(onLaunch, constant.RegisterJobs) {
		for _, j := range jobs {
			if id, err := c.AddFunc("@hourly", j.Handle); err != nil {
				return err
			} else {
				name := j.Name()
				mapping[name] = id
				logger.Infow("Register job", "name", name)
			}
		}
	} else {
		logger.Infof("No [%s] found in [%s]", constant.RegisterJobs, key)
	}
	c.Start()
	return nil
}
