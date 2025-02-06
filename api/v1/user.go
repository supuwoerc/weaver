package v1

import (
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
	if err = u.service.SignUp(ctx, params.ID, params.Code, &models.User{
		Email:    params.Email,
		Password: params.Password,
	}); err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.Success(ctx)
}

func (u *UserApi) Login(ctx *gin.Context) {
	var params request.LoginRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	user, pair, err := u.service.Login(ctx, params.Email, params.Password)
	switch {
	case err == nil && pair != nil:
		response.SuccessWithData(ctx, response.LoginResponse{
			User:         user,
			Token:        pair.AccessToken,
			RefreshToken: pair.RefreshToken,
		})
	default:
		response.FailWithCode(ctx, response.UserLoginFail)
	}
}

func (u *UserApi) Profile(ctx *gin.Context) {
	claims, err := utils.GetContextClaims(ctx)
	if err != nil || claims == nil {
		response.FailWithCode(ctx, response.UserNotExist)
		return
	}
	profile, err := u.service.Profile(ctx, claims.User.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.SuccessWithData(ctx, profile)
}
