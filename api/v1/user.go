package v1

import (
	"errors"
	"gin-web/models"
	"gin-web/pkg/constant"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"gin-web/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"sync"
)

type UserApi struct {
	*BasicApi
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
	service          *service.UserService
}

var (
	userOnce sync.Once
	userApi  *UserApi
)

func NewUserApi() *UserApi {
	userOnce.Do(func() {
		userApi = &UserApi{
			BasicApi:         NewBasicApi(),
			passwordRegexExp: regexp.MustCompile(constant.PasswdRegexPattern, regexp.None),
			service:          service.NewUserService(),
		}
	})
	return userApi
}

// @Tags 用户模块
// @Summary 用户注册
// @Description 用于用户注册帐号
// @Accept json
// @Produce json
// @Param body body request.SignUpRequest true "注册参数"
// @Success 10000 {object} response.BasicResponse[any] "操作成功"
// @Failure 10001 {object} response.BasicResponse[string] "操作失败"
// @Failure 10002 {object} response.BasicResponse[string] "参数错误"
// @Router /api/v1/public/user/signup [post]
func (u *UserApi) SignUp(ctx *gin.Context) {
	var params request.SignUpRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	passwordValid, err := u.passwordRegexExp.MatchString(params.Password)
	if err != nil || !passwordValid {
		response.HttpResponse[any](ctx, response.PasswordValidErr, nil, nil, nil)
		return
	}
	if err = u.service.SignUp(ctx, params.ID, params.Code, models.User{
		Email:    params.Email,
		Password: params.Password,
	}); err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.Success(ctx)
}

// @Tags 用户模块
// @Summary 用户登录
// @Description 用于用户登录
// @Accept json
// @Produce json
// @Param body body request.LoginRequest true "注册参数"
// @Success 10000 {object} response.BasicResponse[any] "操作成功"
// @Failure 10001 {object} response.BasicResponse[any] "操作失败"
// @Failure 10002 {object} response.BasicResponse[any] "参数错误"
// @Router /api/v1/public/user/login [post]
func (u *UserApi) Login(ctx *gin.Context) {
	var params request.LoginRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	user, pair, err := u.service.Login(ctx, params.Email, params.Password)
	switch {
	case pair != nil:
		if err != nil {
			response.FailWithError(ctx, err)
			return
		}
		user.Password = ""
		response.SuccessWithData[response.LoginResponse](ctx, response.LoginResponse{
			User:         user,
			Token:        pair.AccessToken,
			RefreshToken: pair.RefreshToken,
		})
	case errors.Is(err, response.UserLoginFail) || errors.Is(err, response.UserLoginEmailNotFound):
		response.FailWithCode(ctx, response.UserLoginFail)
	default:
		if err != nil {
			response.FailWithMessage(ctx, err.Error())
		}
		response.FailWithError(ctx, err)
	}
}

// @Tags 用户模块
// @Summary 用户信息
// @Description 获取用户账户信息
// @Accept json
// @Produce json
// @Success 10000 {object} response.BasicResponse[any] "操作成功"
// @Failure 10001 {object} response.BasicResponse[any] "操作失败"
// @Failure 10002 {object} response.BasicResponse[any] "参数错误"
// @Router /api/v1/user/profile [get]
func (u *UserApi) Profile(ctx *gin.Context) {
	claims, err := utils.GetContextClaims(ctx)
	if err != nil || claims == nil {
		response.FailWithCode(ctx, response.UserNotExist)
		return
	}
	profile, err := u.service.Profile(ctx, claims.User.UID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	profile.Password = ""
	response.SuccessWithData(ctx, profile)
}

// TODO:补充文档
func (u *UserApi) SetRoles(ctx *gin.Context) {
	var params request.SetRolesRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	err := u.service.SetRoles(ctx, params.UserId, params.RoleIds)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.Success(ctx)
}

// TODO:补充文档
func (u *UserApi) GetRoles(ctx *gin.Context) {
	var params request.GetRolesRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	roles, err := u.service.GetRoles(ctx, params.UserId)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.SuccessWithData(ctx, roles)
}
