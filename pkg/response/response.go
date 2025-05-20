package response

import (
	"context"
	"errors"
	"net/http"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type I18nTranslatorCtxKey string

const (
	I18nTranslatorKey      I18nTranslatorCtxKey = "i18n_translator"
	ValidatorTranslatorKey I18nTranslatorCtxKey = "validator_translator"
)

var (
	namespaceReg = regexp.MustCompile("\\[\\d+\\]", regexp.None)
)

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
	translator, exists := ctx.Get(string(I18nTranslatorKey))
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
				MessageID: code.String(),
			})
		}
	}
	ctx.AbortWithStatusJSON(http.StatusOK, BasicResponse[T]{
		Code:    int(code),
		Data:    data,
		Message: msg,
	})
}

// Success 成功响应-不携带数据
func Success(ctx *gin.Context) {
	HttpResponse[any](ctx, Ok, nil, nil, nil)
}

// SuccessWithData 成功响应-携带数据
func SuccessWithData[T any](ctx *gin.Context, data T) {
	HttpResponse[T](ctx, Ok, data, nil, nil)
}

// SuccessWithMessage 成功响应-携带消息
func SuccessWithMessage(ctx *gin.Context, message string) {
	HttpResponse[any](ctx, Ok, nil, nil, &message)
}

// SuccessWithPageData 成功响应-携带分页数据
func SuccessWithPageData[T any](ctx *gin.Context, total int64, list []T) {
	HttpResponse[DataList[T]](ctx, Ok, DataList[T]{
		Total: total,
		List:  list,
	}, nil, nil)
}

// FailWithMessage 失败响应-携带消息
func FailWithMessage(ctx *gin.Context, message string) {
	HttpResponse[any](ctx, Error, nil, nil, &message)
}

// FailWithCode 失败响应
func FailWithCode(ctx *gin.Context, code StatusCode) {
	HttpResponse[any](ctx, code, nil, nil, nil)
}

// FailWithError 失败响应
func FailWithError(ctx *gin.Context, err error) {
	var code StatusCode
	if errors.As(err, &code) {
		FailWithCode(ctx, code)
		return
	}
	if errors.Is(err, context.Canceled) {
		FailWithCode(ctx, CancelRequest)
		return
	}
	if errors.Is(err, context.DeadlineExceeded) {
		FailWithCode(ctx, TimeoutErr)
		return
	}
	FailWithMessage(ctx, err.Error())
}

// ParamsValidateFail 失败响应-参数错误
func ParamsValidateFail(ctx *gin.Context, err error) {
	msg := err.Error()
	var errs validator.ValidationErrors
	ok := errors.As(err, &errs)
	if !ok {
		HttpResponse[any](ctx, InvalidParams, nil, nil, &msg)
	} else if translator, exists := ctx.Get(string(ValidatorTranslatorKey)); exists {
		if trans, isOk := translator.(ut.Translator); isOk {
			errMap := make(map[string]string)
			for _, e := range errs {
				fieldName := e.Field()
				replace, temp := namespaceReg.Replace(fieldName, "", 0, -1)
				if temp == nil {
					errMap[replace] = e.Translate(trans)
				} else {
					errMap[fieldName] = e.Translate(trans)
				}
			}
			HttpResponse[any](ctx, InvalidParams, errMap, nil, nil)
		} else {
			HttpResponse[any](ctx, InvalidParams, nil, nil, &msg)
		}
	} else {
		HttpResponse[any](ctx, InvalidParams, nil, nil, &msg)
	}
}
