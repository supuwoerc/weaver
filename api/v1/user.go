package v1

import (
	"context"
	"errors"
	"gin-web/models"
	"gin-web/pkg/constant"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"gin-web/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

type UserService interface {
	SignUp(ctx context.Context, id string, code string, user *models.User) error
	Login(ctx context.Context, email string, password string) (*response.LoginResponse, error)
	Profile(ctx context.Context, uid uint) (*response.ProfileResponse, error)
	GetUserList(ctx context.Context, keyword string, limit, offset int) ([]*response.UserListRowResponse, int64, error)
}

type UserApi struct {
	*BasicApi
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
	phoneRegexExp    *regexp.Regexp
	service          UserService
}

var (
	userOnce sync.Once
	userApi  *UserApi
)

func NewUserApi() *UserApi {
	userOnce.Do(func() {
		userApi = &UserApi{
			BasicApi:         NewBasicApi(),
			emailRegexExp:    regexp.MustCompile(constant.EmailRegexPattern, regexp.None),
			passwordRegexExp: regexp.MustCompile(constant.PasswdRegexPattern, regexp.None),
			phoneRegexExp:    regexp.MustCompile(constant.PhoneRegexPattern, regexp.None),
			service:          service.NewUserService(),
		}
	})
	return userApi
}

func (r *UserApi) SignUp(ctx *gin.Context) {
	var params request.SignUpRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	emailValid, err := r.emailRegexExp.MatchString(params.Email)
	if err != nil || !emailValid {
		response.HttpResponse[any](ctx, response.EmailValidErr, nil, nil, nil)
		return
	}
	passwordValid, err := r.passwordRegexExp.MatchString(params.Password)
	if err != nil || !passwordValid {
		response.HttpResponse[any](ctx, response.PasswordValidErr, nil, nil, nil)
		return
	}
	err = r.service.SignUp(ctx, params.ID, params.Code, &models.User{
		Email:    params.Email,
		Password: params.Password,
	})
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.Success(ctx)
}

func (r *UserApi) Login(ctx *gin.Context) {
	var params request.LoginRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	res, err := r.service.Login(ctx, params.Email, params.Password)
	if err != nil {
		if errors.Is(err, response.UserInactive) || errors.Is(err, response.UserDisabled) {
			response.FailWithError(ctx, err)
		} else {
			response.FailWithCode(ctx, response.UserLoginFail)
		}
		return
	}
	response.SuccessWithData(ctx, res)
}

func (r *UserApi) Profile(ctx *gin.Context) {
	claims, err := utils.GetContextClaims(ctx)
	if err != nil || claims == nil {
		response.FailWithCode(ctx, response.UserNotExist)
		return
	}
	detail, err := r.service.Profile(ctx, claims.User.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.SuccessWithData(ctx, detail)
}

func (r *UserApi) GetUserList(ctx *gin.Context) {
	var params request.GetUserListRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	list, total, err := r.service.GetUserList(ctx, params.Keyword, params.Limit, params.Offset)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.SuccessWithPageData(ctx, total, list)
}

func (r *UserApi) ActiveSuccess(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "active-success.html", nil)
}

func (r *UserApi) ActiveFailure(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "active-failure.html", nil)
}
