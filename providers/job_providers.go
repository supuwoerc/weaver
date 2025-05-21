package providers

import (
	"time"

	"github.com/google/wire"
	"github.com/supuwoerc/weaver/pkg/job"
	"github.com/supuwoerc/weaver/pkg/logger"
)

func SystemJobs(logger *logger.Logger) []job.SystemJob {
	return []job.SystemJob{
		job.NewServerStatus(10*time.Second, logger),
	}
}

var SystemJobProvider = wire.NewSet(
	SystemJobs,
	job.NewSystemJobManager,
)
