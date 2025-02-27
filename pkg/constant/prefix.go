package constant

type Prefix string

const (
	PermissionIdPrefix       Prefix = "lock:permission:id"
	PermissionNamePrefix     Prefix = "lock:permission:name"
	PermissionResourcePrefix Prefix = "lock:permission:resource"
	RoleIdPrefix             Prefix = "lock:role:id"
	RoleNamePrefix           Prefix = "lock:role:name"
	SignUpEmailPrefix        Prefix = "lock:signup:email"
	DepartmentIdPrefix       Prefix = "lock:department:id"
	DepartmentNamePrefix     Prefix = "lock:department:name"
)

const (
	CaptchaCodePrefix   Prefix = "captcha:"
	ActiveAccountPrefix Prefix = "active:account:"
)
