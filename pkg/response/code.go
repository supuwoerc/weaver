package response

type StatusCode = int

// 响应的code枚举
const (
	SUCCESS                   StatusCode = 10000 // 通用成功
	ERROR                                = 10001 // 通用错误
	INVALID_PARAMS                       = 10002 // 错误的参数
	INVALID_TOKEN                        = 10003 // token错误
	UNKNOWN_ERROR                        = 10004 // 未知错误
	RECOVERY_ERROR                       = 10005 // 发生recovery
	INVALID_REFRESH_TOKEN                = 10006 // 长token错误
	UNNECESSARY_REFRESH_TOKEN            = 10007 // 不必要的刷新token(短token还未过期)
	CASBIN_ERR                           = 10008 // casbin校验出错
	CASBIN_INVALID                       = 10009 // casbin校验未通过
)

const (
	USER_CREATE_DUPLICATE_EMAIL     StatusCode = 20000 // 用户创建email重复导致的唯一索引错误
	USER_LOGIN_EMAIL_NOT_FOUND                 = 20001 // 用户登录email未查询到
	USER_LOGIN_FAIL                            = 20002 // 用户登录失败
	USER_LOGIN_TOKEN_PAIR_CACHE_ERR            = 20003 // 用户登录时的redis存储token对失败
	PASSWORD_VALID_ERR                         = 20004 // 密码格式错误
)

const (
	CAPTCHA_VERIFY_FAIL StatusCode = 30000 // 验证码校验错误
)

const (
	ROLE_CREATE_DUPLICATE_NAME StatusCode = 40000 // 角色名已存在
)
