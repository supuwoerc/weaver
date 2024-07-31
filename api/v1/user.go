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
	service          func(ctx *gin.Context) *service.UserService
}

func NewUserApi() UserApi {
	return UserApi{
		BasicApi:         NewBasicApi(),
		passwordRegexExp: regexp.MustCompile(constant.PasswdRegexPattern, regexp.None),
		service: func(ctx *gin.Context) *service.UserService {
			return service.NewUserService(ctx)
		},
	}
}

// @Tags 用户管理模块
// @Summary 用户注册
// @Description 用于用户注册帐号
// @Accept json
// @Produce json
// @Param body body request.SignUpRequest true "注册参数"
// @Success 10000 {object} response.BasicResponse[any] "操作成功"
// @Failure 10001 {object} response.BasicResponse[any] "操作失败"
// @Failure 10002 {object} response.BasicResponse[any] "参数错误"
// @Router /api/v1/public/user/signup [post]
func (u UserApi) SignUp(ctx *gin.Context) {
	var params request.SignUpRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	passwordValid, err := u.passwordRegexExp.MatchString(params.Password)
	if err != nil || !passwordValid {
		response.HttpResponse[any](ctx, response.PASSWORD_VALID_ERR, nil, nil, nil)
		return
	}
	if err = u.service(ctx).SignUp(ctx.Request.Context(), params.ID, params.Code, models.User{
		Email:    params.Email,
		Password: &params.Password,
	}); err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.Success(ctx)
}

// @Tags 用户管理模块
// @Summary 用户登录
// @Description 用于用户登录
// @Accept json
// @Produce json
// @Param body body request.LoginRequest true "注册参数"
// @Success 10000 {object} response.BasicResponse[any] "操作成功"
// @Failure 10001 {object} response.BasicResponse[any] "操作失败"
// @Failure 10002 {object} response.BasicResponse[any] "参数错误"
// @Router /api/v1/public/user/login [post]
func (u UserApi) Login(ctx *gin.Context) {
	var params request.LoginRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	user, pair, err := u.service(ctx).Login(ctx.Request.Context(), params.Email, params.Password)
	switch {
	case pair != nil:
		user.Password = nil
		if err != nil {
			response.FailWithError(ctx, err)
			return
		}
		response.SuccessWithData[response.LoginResponse](ctx, response.LoginResponse{
			User:         user,
			Token:        pair.AccessToken,
			RefreshToken: pair.RefreshToken,
		})
	case err == constant.GetError(ctx, response.USER_LOGIN_FAIL) || err == constant.GetError(ctx, response.USER_LOGIN_EMAIL_NOT_FOUND):
		response.FailWithCode(ctx, response.USER_LOGIN_FAIL)
	default:
		response.FailWithMessage(ctx, err.Error())
	}
}

// 获取个人信息
func (u UserApi) Profile(ctx *gin.Context) {
	response.Success(ctx)
}
