package constant

//go:generate stringer -type=JobType -linecomment -output job_string.go
type JobType int

const (
	ServerStatus JobType = iota + 1 // serverStatus
)
