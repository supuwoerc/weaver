package response

import "fmt"

//go:generate stringer -type=StatusCode -linecomment -output code_string.go
type StatusCode int

func (s StatusCode) Error() string {
	return fmt.Sprintf("%d", s)
}

// 响应的code枚举
const (
	Ok                      StatusCode = 10000 // ok
	Error                   StatusCode = 10001 // error
	InvalidParams           StatusCode = 10002 // invalidParams
	InvalidToken            StatusCode = 10003 // invalidToken
	CancelRequest           StatusCode = 10004 // cancelRequest
	RecoveryError           StatusCode = 10005 // recoveryError
	InvalidRefreshToken     StatusCode = 10006 // invalidRefreshToken
	UnnecessaryRefreshToken StatusCode = 10007 // unnecessaryRefreshToken
	AuthErr                 StatusCode = 10008 // authErr
	NoAuthority             StatusCode = 10009 // noAuthority
	TimeoutErr              StatusCode = 10010 // timeoutErr
	Busy                    StatusCode = 10011 // busy
)

const (
	UserCreateDuplicateEmail   StatusCode = 20000 // userCreateDuplicateEmail
	UserLoginFail              StatusCode = 20001 // userLoginFail
	UserLoginTokenPairCacheErr StatusCode = 20002 // userLoginTokenPairCacheErr
	PasswordValidErr           StatusCode = 20003 // passwordValidErr
	UserNotExist               StatusCode = 20004 // userNotExist
)

const (
	CaptchaVerifyFail StatusCode = 30000 // captchaVerifyFail
)

const (
	RoleCreateDuplicateName StatusCode = 40000 // roleCreateDuplicateName
	NoValidRoles            StatusCode = 40001 // noValidRoles
	RoleNotExist            StatusCode = 40002 // roleNotExist
	RoleExistPermissionRef  StatusCode = 40003 // roleExistPermissionRef
	RoleExistUserRef        StatusCode = 40004 // roleExistUserRef
)

const (
	PermissionCreateDuplicate StatusCode = 50000 // permissionCreateDuplicate
	PermissionNotExist        StatusCode = 50001 // permissionNotExist
	PermissionExistRoleRef    StatusCode = 50002 // permissionExistRoleRef
)

const (
	InvalidAttachmentLength StatusCode = 50000 // invalidAttachmentLength
)

const (
	DeptCreateDuplicate StatusCode = 60000 // deptCreateDuplicate
	DeptNotExist        StatusCode = 60001 // deptNotExist
	DeptExistUserRef    StatusCode = 60002 // deptExistUserRef
)
