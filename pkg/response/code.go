package response

type StatusCode = int

// 响应的code枚举
const (
	Ok                      StatusCode = 10000 // 通用成功
	Error                              = 10001 // 通用错误
	InvalidParams                      = 10002 // 错误的参数
	InvalidToken                       = 10003 // token错误
	CancelRequest                      = 10004 // 请求取消
	RecoveryError                      = 10005 // 发生recovery
	InvalidRefreshToken                = 10006 // 长token错误
	UnnecessaryRefreshToken            = 10007 // 不必要的刷新token(短token还未过期)
	CasbinErr                          = 10008 // casbin校验出错
	CasbinInvalid                      = 10009 // casbin校验未通过
	TimeoutErr                         = 10010 // 上下文超时取消(context.WithTimeout)
)

const (
	UserCreateDuplicateEmail   StatusCode = 20000 // 用户创建email重复导致的唯一索引错误
	UserLoginEmailNotFound                = 20001 // 用户登录email未查询到
	UserLoginFail                         = 20002 // 用户登录失败
	UserLoginTokenPairCacheErr            = 20003 // 用户登录时的redis存储token对失败
	PasswordValidErr                      = 20004 // 密码格式错误
	UserNotExist                          = 20005 // 用户不存在(管理员设置用户相关内容时)
)

const (
	CaptchaVerifyFail StatusCode = 30000 // 验证码校验错误
)

const (
	RoleCreateDuplicateName StatusCode = 40000 // 角色名已存在
	NoValidRoles            StatusCode = 40001 // 设置角色时无有效角色(角色ID全部无效)
)

const (
	InvalidAttachmentLength StatusCode = 50000 // 文件上传数量不合法
)
