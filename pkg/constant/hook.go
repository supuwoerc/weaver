package constant

//go:generate stringer -type=HookType -linecomment -output hook_string.go
type HookType int

const (
	RegisterJobs HookType = iota + 1 // registerJobs
)
