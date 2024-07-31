package response

import (
	"errors"
	"gin-web/pkg/global"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"net/http"
	"strconv"
)

const TranslatorKey = "translator"
const Locale = "locale"

// BasicResponse 通用的数据返回
type BasicResponse[T any] struct {
	Code    int    `json:"code"`
	Data    T      `json:"data"`
	Message string `json:"message"`
}

// DataList 列表类型的数据
type DataList[T any] struct {
	Total int64 `json:"total"`
	List  []T   `json:"list"`
}

// HttpResponse json响应
func HttpResponse[T any](ctx *gin.Context, code StatusCode, data T, config *i18n.LocalizeConfig, message *string) {
	translator, exists := ctx.Get(TranslatorKey)
	var msg string
	if message != nil {
		msg = *message
	}
	if exists && message == nil {
		loc := translator.(*i18n.Localizer)
		if config != nil {
			msg = loc.MustLocalize(config)
		} else {
			msg = loc.MustLocalize(&i18n.LocalizeConfig{
				MessageID: strconv.Itoa(code),
			})
		}
	}
	ctx.AbortWithStatusJSON(http.StatusOK, BasicResponse[T]{
		Code:    code,
		Data:    data,
		Message: msg,
	})
}

// Success 成功响应-不携带数据
func Success(ctx *gin.Context) {
	HttpResponse[any](ctx, SUCCESS, nil, nil, nil)
}

// SuccessWithData 成功响应-携带数据
func SuccessWithData[T any](ctx *gin.Context, data T) {
	HttpResponse[T](ctx, SUCCESS, data, nil, nil)
}

// SuccessWithMessage 成功响应-携带消息
func SuccessWithMessage(ctx *gin.Context, message string) {
	HttpResponse[any](ctx, SUCCESS, nil, nil, &message)
}

// SuccessWithPageData 成功响应-携带分页数据
func SuccessWithPageData[T any](ctx *gin.Context, total int64, list []T) {
	HttpResponse[DataList[T]](ctx, SUCCESS, DataList[T]{
		Total: total,
		List:  list,
	}, nil, nil)
}

// FailWithMessage 失败响应-携带消息
func FailWithMessage(ctx *gin.Context, message string) {
	HttpResponse[any](ctx, ERROR, nil, nil, &message)
}

// FailWithCode 失败响应
func FailWithCode(ctx *gin.Context, code StatusCode) {
	HttpResponse[any](ctx, code, nil, nil, nil)
}

// FailWithError 失败响应
func FailWithError(ctx *gin.Context, err error) {
	for code, e := range global.LocaleErrors[global.CN] {
		if errors.Is(e, err) {
			FailWithCode(ctx, code)
			return
		}
	}
	for code, e := range global.LocaleErrors[global.EN] {
		if errors.Is(e, err) {
			FailWithCode(ctx, code)
			return
		}
	}
	FailWithMessage(ctx, err.Error())
}

// ParamsValidateFail 失败响应-参数错误
func ParamsValidateFail(ctx *gin.Context, err error) {
	msg := err.Error()
	HttpResponse[any](ctx, INVALID_PARAMS, nil, nil, &msg)
}
