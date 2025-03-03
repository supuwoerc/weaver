package job

import "gin-web/pkg/constant"

type SystemJob interface {
	Name() string
	IfStillRunning() constant.JobStillMode
	Handle()
	Interval() string
}
