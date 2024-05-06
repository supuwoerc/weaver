package v1

import (
	"gin-web/models"
	"gin-web/pkg/constant"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
)

type UserApi struct {
	*BasicApi
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
	service          *service.UserService
}

func NewUserApi() UserApi {
	return UserApi{
		BasicApi:         NewBasicApi(),
		passwordRegexExp: regexp.MustCompile(constant.PasswdRegexPattern, regexp.None),
		service:          service.NewUserService(),
	}
}

// 用户注册
func (u UserApi) SignUp(ctx *gin.Context) {
	var params request.SignUpRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ParamsValidateFail(ctx)
		return
	}
	passwordValid, err := u.passwordRegexExp.MatchString(params.Password)
	if err != nil || !passwordValid {
		response.FailWithMessage(ctx, "密码格式错误")
		return
	}
	if err = u.service.SignUp(ctx.Request.Context(), models.User{
		Email:    params.Email,
		Password: params.Password,
	}); err != nil {
		response.FailWithMessage(ctx, err.Error())
		return
	}
	response.Success(ctx)
}
