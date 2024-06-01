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

// NOTE:USER模块-START

// 创建用户时数据库唯一索引错误
var USER_CREATE_DUPLICATE_EMAIL_ERR = errors.New(response.GetMessage(response.USER_CREATE_DUPLICATE_EMAIL))

// 用户登录email未查询到错误
var USER_LOGIN_EMAIL_NOT_FOUND_ERR = errors.New(response.GetMessage(response.USER_LOGIN_EMAIL_NOT_FOUND))

// 用户登录失败
var USER_LOGIN_FAIL_ERR = errors.New(response.GetMessage(response.USER_LOGIN_FAIL))

// 用户登录时的redis存储token对失败
var USER_LOGIN_TOKEN_PAIR_CACHE_ERR = errors.New(response.GetMessage(response.USER_LOGIN_TOKEN_PAIR_CACHE_ERR))

// NOTE:USER模块-END
