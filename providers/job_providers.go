package providers

import (
	"time"

	"github.com/supuwoerc/weaver/pkg/job"

	"github.com/google/wire"
	"go.uber.org/zap"
)

func SystemJobs(logger *zap.SugaredLogger) []job.SystemJob {
	return []job.SystemJob{
		job.NewServerStatus(10*time.Second, logger),
	}
}

var SystemJobProvider = wire.NewSet(
	SystemJobs,
	job.NewSystemJobManager,
)
