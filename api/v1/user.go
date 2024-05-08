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

// @Tags 用户管理模块
// @Summary 用户登录
// @Description 用于用户注册帐号
// @Accept json
// @Produce json
// @Param body body request.SignUpRequest true "注册参数"
// @Success 10000 {object} response.BasicResponse[any] "操作成功"
// @Failure 10001 {object} response.BasicResponse[any] "操作失败"
// @Failure 10002 {object} response.BasicResponse[any] "参数错误"
// @Failure 20000 {object} response.BasicResponse[any] "邮箱已注册"
// @Router /api/v1/public/user/signup [post]
func (u UserApi) SignUp(ctx *gin.Context) {
	// {"code":10000,"message":"操作成功"
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
