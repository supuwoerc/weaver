package response

// 响应的code枚举
const (
	SUCCESS                   int = 10000 // 通用成功
	ERROR                         = 10001 // 通用错误
	INVALID_PARAMS                = 10002 // 错误的参数
	INVALID_TOKEN                 = 10003 // 短token失效
	UNKNOWN_ERROR                 = 10004 // 未知错误
	RECOVERY_ERROR                = 10005 // 发生recovery
	INVALID_REFRESH_TOKEN         = 10006 // 长token失效
	UNNECESSARY_REFRESH_TOKEN     = 1007  // 不必要的刷新token(短token还未过期)
)
