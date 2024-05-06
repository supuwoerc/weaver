package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 通用的数据返回
type BasicResponse[T any] struct {
	Code    int    `json:"code"`
	Data    T      `json:"data"`
	Message string `json:"message"`
}

// 列表类型的数据
type DataList[T any] struct {
	Total int64 `json:"total"`
	List  []T   `json:"list"`
}

// json响应
func HttpResponse[T any](ctx *gin.Context, code int, data T, message string) {
	msg := GetMessage(code)
	if message != "" {
		msg = message
	}
	ctx.AbortWithStatusJSON(http.StatusOK, BasicResponse[T]{
		Code:    code,
		Data:    data,
		Message: msg,
	})
}

// 成功响应-不携带数据
func Success(ctx *gin.Context) {
	HttpResponse[any](ctx, SUCCESS, nil, "")
}

// 成功响应-携带数据
func SuccessWithData[T any](ctx *gin.Context, data T) {
	HttpResponse[T](ctx, SUCCESS, data, "")
}

// 成功响应-携带消息
func SuccessWithMessage(ctx *gin.Context, message string) {
	HttpResponse[any](ctx, SUCCESS, nil, message)
}

// 成功响应-携带分页数据
func SuccessWithPageData[T any](ctx *gin.Context, total int64, list []T) {
	HttpResponse[DataList[T]](ctx, SUCCESS, DataList[T]{
		Total: total,
		List:  list,
	}, "")
}

// 失败响应-不携带数据
func Fail(ctx *gin.Context) {
	HttpResponse[any](ctx, ERROR, nil, "")
}

// 失败响应-携带数据
func FailWithData[T any](ctx *gin.Context, data T) {
	HttpResponse[any](ctx, ERROR, data, "")
}

// 失败响应-携带消息
func FailWithMessage(ctx *gin.Context, message string) {
	HttpResponse[any](ctx, ERROR, nil, message)
}
