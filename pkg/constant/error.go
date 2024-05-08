package constant

import (
	"errors"
	"gin-web/pkg/response"
)

// 短TOKEN解析错误
var TOKEN_PARSE_ERROR = errors.New(response.GetMessage(response.INVALID_TOKEN))

// 长TOKEN解析错误
var REFRESH_TOKEN_PARSE_ERROR = errors.New(response.GetMessage(response.INVALID_REFRESH_TOKEN))

// 不必要的刷新短TOKEN
var UNNECESSARY_REFRESH_TOKEN_ERROR = errors.New(response.GetMessage(response.UNNECESSARY_REFRESH_TOKEN))

// 数据库唯一索引错误
var USER_CREATE_DUPLICATE_EMAIL_ERR = errors.New(response.GetMessage(response.USER_CREATE_DUPLICATE_EMAIL))
