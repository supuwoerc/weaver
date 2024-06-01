package response

// 响应的code枚举
const (
	SUCCESS                   int = 10000 // 通用成功
	ERROR                         = 10001 // 通用错误
	INVALID_PARAMS                = 10002 // 错误的参数
	INVALID_TOKEN                 = 10003 // token错误
	UNKNOWN_ERROR                 = 10004 // 未知错误
	RECOVERY_ERROR                = 10005 // 发生recovery
	INVALID_REFRESH_TOKEN         = 10006 // 长token错误
	UNNECESSARY_REFRESH_TOKEN     = 10007 // 不必要的刷新token(短token还未过期)
)

const (
	USER_CREATE_DUPLICATE_EMAIL     int = 20000 // 用户创建email重复导致的唯一索引错误
	USER_LOGIN_EMAIL_NOT_FOUND      int = 20001 // 用户登录email未查询到
	USER_LOGIN_FAIL                 int = 20002 // 用户登录失败
	USER_LOGIN_TOKEN_PAIR_CACHE_ERR int = 20003 // 用户登录时的redis存储token对失败
)
