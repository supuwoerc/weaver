package bootstrap

import (
	"gin-web/initialize"
	"gin-web/pkg/constant"
	"gin-web/pkg/global"
	"gin-web/pkg/job"
	"github.com/robfig/cron/v3"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"time"
)

var (
	jobs    []job.SystemJob
	mapping = make(map[string]cron.EntryID)
)

func init() {
	jobs = []job.SystemJob{
		job.NewServerStatus(12 * time.Second),
	}
}

func skip(f func()) cron.Job {
	w := cron.FuncJob(f)
	// https://github.com/robfig/cron/issues/366
	wrapJob := cron.NewChain(cron.SkipIfStillRunning(global.CronLogger), initialize.CronRecover(global.CronLogger)).Then(&w)
	return wrapJob
}

func delay(f func()) cron.Job {
	w := cron.FuncJob(f)
	wrapJob := cron.NewChain(cron.DelayIfStillRunning(global.CronLogger), initialize.CronRecover(global.CronLogger)).Then(&w)
	return wrapJob
}

func RegisterJobs() error {
	key := "system.hooks.launch"
	onLaunch := viper.GetStringSlice(key)
	if lo.Contains(onLaunch, constant.RegisterJobs) {
		for _, j := range jobs {
			mode := j.IfStillRunning()
			var id cron.EntryID
			var err error
			switch mode {
			case constant.Skip:
				id, err = global.Cron.AddJob("@hourly", skip(j.Handle))
			case constant.Delay:
				id, err = global.Cron.AddJob("@hourly", delay(j.Handle))
			case constant.None:
				fallthrough
			default:
				id, err = global.Cron.AddFunc("@hourly", j.Handle)
			}
			if err != nil {
				return err
			}
			name := j.Name()
			mapping[name] = id
			global.Logger.Infow("Register job", "name", name)
		}
	} else {
		global.Logger.Infof("No [%s] found in [%s]", constant.RegisterJobs, key)
	}
	global.Cron.Start()
	return nil
}
