package constant

//go:generate stringer -type=Prefix -linecomment -output prefix_string.go
type Prefix int

const (
	PermissionIdPrefix       Prefix = iota + 1 // lock:permission:id
	PermissionNamePrefix                       // lock:permission:name
	PermissionResourcePrefix                   // lock:permission:resource
	RoleIdPrefix                               // lock:role:id
	RoleNamePrefix                             // lock:role:name
	SignUpEmailPrefix                          // lock:signup:email
)
