package constant

import (
	"errors"
	"gin-web/pkg/global"
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/viper"
	"strconv"
)

type StatusCode2Error map[response.StatusCode]error

var cnErrorMap StatusCode2Error
var enErrorMap StatusCode2Error

var systemModuleCode = []response.StatusCode{
	response.INVALID_TOKEN,
	response.CANCEL_REQUEST,
	response.TIMEOUT_ERROR,
	response.INVALID_REFRESH_TOKEN,
	response.UNNECESSARY_REFRESH_TOKEN,
	response.CASBIN_ERR,
	response.CASBIN_INVALID,
}
var userModuleCode = []response.StatusCode{
	response.USER_CREATE_DUPLICATE_EMAIL,
	response.USER_LOGIN_EMAIL_NOT_FOUND,
	response.USER_LOGIN_FAIL,
	response.USER_LOGIN_TOKEN_PAIR_CACHE_ERR,
	response.USER_NOT_EXIST,
}
var captchaModuleCode = []response.StatusCode{
	response.CAPTCHA_VERIFY_FAIL,
}
var roleModuleCode = []response.StatusCode{
	response.ROLE_CREATE_DUPLICATE_NAME,
	response.NO_VALID_ROLES,
}

var codeModules = [][]response.StatusCode{systemModuleCode, userModuleCode, captchaModuleCode, roleModuleCode}

type InitParams struct {
	CN *i18n.Localizer
	EN *i18n.Localizer
}

func initWithLang(localizer *i18n.Localizer, codes []response.StatusCode, sourceMap *StatusCode2Error) {
	for _, code := range codes {
		msg := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: strconv.Itoa(code),
		})
		(*sourceMap)[code] = errors.New(msg)
	}
}

func InitErrors(localizer InitParams) map[string]map[int]error {
	var codes []response.StatusCode
	if cnErrorMap == nil || enErrorMap == nil {
		for _, module := range codeModules {
			codes = append(codes, module...)
		}
	}
	if cnErrorMap == nil {
		cnErrorMap = StatusCode2Error{}
		initWithLang(localizer.CN, codes, &cnErrorMap)
	}
	if enErrorMap == nil {
		enErrorMap = StatusCode2Error{}
		initWithLang(localizer.EN, codes, &enErrorMap)
	}
	return map[string]map[int]error{
		"cn": cnErrorMap,
		"en": enErrorMap,
	}
}

func GetError(ctx *gin.Context, code response.StatusCode) error {
	value, exists := ctx.Get(response.Locale)
	locale := ""
	if exists {
		locale = value.(string)
	} else {
		locale = viper.GetString("system.defaultLang")
	}
	if locale == global.CN {
		return cnErrorMap[code]
	}
	return enErrorMap[code]
}
