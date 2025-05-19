package job

import (
	"gin-web/initialize"
	"gin-web/pkg/constant"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type SystemJob interface {
	Name() string
	IfStillRunning() constant.JobStillMode
	Handle()
	Interval() string
}

type SystemJobManager struct {
	cronLogger *initialize.CronLogger
	cronClient *cron.Cron
	logger     *zap.SugaredLogger
	jobsMap    map[string]cron.EntryID // TODO:动态开关任务
	jobs       []SystemJob             // 任务集合
}

func NewSystemJobManager(cl *initialize.CronLogger, c *cron.Cron, logger *zap.SugaredLogger, j ...SystemJob) *SystemJobManager {
	return &SystemJobManager{
		cronLogger: cl,
		cronClient: c,
		logger:     logger,
		jobsMap:    make(map[string]cron.EntryID),
		jobs:       j,
	}
}

func (jr *SystemJobManager) skip(f func()) cron.Job {
	w := cron.FuncJob(f)
	// https://github.com/robfig/cron/issues/366
	wrapJob := cron.NewChain(cron.SkipIfStillRunning(jr.cronLogger), jr.cronLogger.CronRecover()).Then(&w)
	return wrapJob
}

func (jr *SystemJobManager) delay(f func()) cron.Job {
	w := cron.FuncJob(f)
	wrapJob := cron.NewChain(cron.DelayIfStillRunning(jr.cronLogger), jr.cronLogger.CronRecover()).Then(&w)
	return wrapJob
}

func (jr *SystemJobManager) RegisterJobsAndStart() error {
	for _, j := range jr.jobs {
		var id cron.EntryID
		var err error
		mode := j.IfStillRunning()
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
		jr.jobsMap[name] = id
		jr.logger.Infow("Register job", "name", name, "interval", j.Interval(), "id", id)
	}
	jr.cronClient.Start()
	return nil
}

func (jr *SystemJobManager) Stop() {
	ctx := jr.cronClient.Stop()
	<-ctx.Done()
	jr.logger.Info("SystemJobManager:cron jobs have been stopped")
}
