package job

import (
	"gin-web/initialize"
	"gin-web/pkg/constant"
	"github.com/robfig/cron/v3"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"sync"
	"time"
)

type SystemJob interface {
	Name() string
	IfStillRunning() constant.JobStillMode
	Handle()
	Interval() string
}

var (
	mapping = make(map[string]cron.EntryID)
)

type JobRegisterer struct {
	cronLogger *initialize.CronLogger
	cronClient *cron.Cron
	logger     *zap.SugaredLogger
	viper      *viper.Viper
}

func NewJobRegisterer(cl *initialize.CronLogger, c *cron.Cron, logger *zap.SugaredLogger, v *viper.Viper) *JobRegisterer {
	return &JobRegisterer{
		cronLogger: cl,
		cronClient: c,
		logger:     logger,
		viper:      v,
	}
}

func (jr *JobRegisterer) skip(f func()) cron.Job {
	w := cron.FuncJob(f)
	// https://github.com/robfig/cron/issues/366
	wrapJob := cron.NewChain(cron.SkipIfStillRunning(jr.cronLogger), jr.cronLogger.CronRecover()).Then(&w)
	return wrapJob
}

func (jr *JobRegisterer) delay(f func()) cron.Job {
	w := cron.FuncJob(f)
	wrapJob := cron.NewChain(cron.DelayIfStillRunning(jr.cronLogger), jr.cronLogger.CronRecover()).Then(&w)
	return wrapJob
}

func (jr *JobRegisterer) initSystemJobs() []SystemJob {
	return []SystemJob{
		NewServerStatus(10*time.Second, jr.logger),
	}
}

func (jr *JobRegisterer) RegisterJobsAndStart() error {
	key := "system.hooks.launch"
	onLaunch := jr.viper.GetStringSlice(key)
	if lo.Contains(onLaunch, constant.RegisterJobs) {
		for _, j := range jr.initSystemJobs() {
			mode := j.IfStillRunning()
			var id cron.EntryID
			var err error
			switch mode {
			case constant.Skip:
				id, err = jr.cronClient.AddJob(j.Interval(), jr.skip(j.Handle))
			case constant.Delay:
				id, err = jr.cronClient.AddJob(j.Interval(), jr.delay(j.Handle))
			case constant.None:
				fallthrough
			default:
				id, err = jr.cronClient.AddFunc(j.Interval(), j.Handle)
			}
			if err != nil {
				return err
			}
			name := j.Name()
			mapping[name] = id
			jr.logger.Infow("Register job", "name", name, "interval", j.Interval(), "id", id)
		}
	} else {
		jr.logger.Infof("No [%s] found in [%s]", constant.RegisterJobs, key)
	}
	jr.cronClient.Start()
	return nil
}

func (jr *JobRegisterer) Stop(group *sync.WaitGroup) {
	defer group.Done()
	ctx := jr.cronClient.Stop()
	<-ctx.Done()
	jr.logger.Info("JobRegisterer:cron jobs have been stopped")
}
