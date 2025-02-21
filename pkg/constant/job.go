package constant

// JobName 定时任务名称
type JobName string

const (
	ServerStatus JobName = "serverStatus"
)

// JobStillMode 上一个定时任务还在执行中,当前任务的模式
type JobStillMode int

const (
	Skip JobStillMode = iota + 1
	Delay
	None
)
