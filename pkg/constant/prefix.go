package constant

//go:generate stringer -type=Prefix -linecomment -output prefix_string.go
type Prefix int

const (
	PermissionIdLockPrefix Prefix = iota + 1 // lock:permission
	RoleIdLockPrefix                         // lock:role
)
