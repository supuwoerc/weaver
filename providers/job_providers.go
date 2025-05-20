package providers

import (
	"time"

	"github.com/supuwoerc/weaver/pkg/job"
	"github.com/supuwoerc/weaver/pkg/logger"

	"github.com/google/wire"
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
