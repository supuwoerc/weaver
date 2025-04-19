package job

import (
	"gin-web/conf"
	"gin-web/initialize"
	"gin-web/pkg/constant"
	"github.com/robfig/cron/v3"
	"github.com/samber/lo"
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

type SystemJobRegisterer struct {
	cronLogger *initialize.CronLogger
	cronClient *cron.Cron
	logger     *zap.SugaredLogger
	conf       *conf.Config
}

func NewJobRegisterer(cl *initialize.CronLogger, c *cron.Cron, logger *zap.SugaredLogger, conf *conf.Config) *SystemJobRegisterer {
	return &SystemJobRegisterer{
		cronLogger: cl,
		cronClient: c,
		logger:     logger,
		conf:       conf,
	}
}

func (jr *SystemJobRegisterer) skip(f func()) cron.Job {
	w := cron.FuncJob(f)
	// https://github.com/robfig/cron/issues/366
	wrapJob := cron.NewChain(cron.SkipIfStillRunning(jr.cronLogger), jr.cronLogger.CronRecover()).Then(&w)
	return wrapJob
}

func (jr *SystemJobRegisterer) delay(f func()) cron.Job {
	w := cron.FuncJob(f)
	wrapJob := cron.NewChain(cron.DelayIfStillRunning(jr.cronLogger), jr.cronLogger.CronRecover()).Then(&w)
	return wrapJob
}

func (jr *SystemJobRegisterer) initSystemJobs() []SystemJob {
	return []SystemJob{
		NewServerStatus(10*time.Second, jr.logger),
	}
}

func (jr *SystemJobRegisterer) RegisterJobsAndStart() error {
	onLaunch := jr.conf.System.Hooks.Launch
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
		jr.logger.Infof("No [%s] found in hooks config", constant.RegisterJobs)
	}
	jr.cronClient.Start()
	return nil
}

func (jr *SystemJobRegisterer) Stop(group *sync.WaitGroup) {
	defer group.Done()
	ctx := jr.cronClient.Stop()
	<-ctx.Done()
	jr.logger.Info("SystemJobRegisterer:cron jobs have been stopped")
}
