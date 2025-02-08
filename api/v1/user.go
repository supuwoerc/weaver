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
	"github.com/samber/lo"
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
	err = u.service.SignUp(ctx, params.ID, params.Code, &models.User{
		Email:    params.Email,
		Password: params.Password,
	})
	if err != nil {
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
	if err != nil {
		response.FailWithCode(ctx, response.UserLoginFail)
		return
	}
	response.SuccessWithData(ctx, &response.LoginResponse{
		User: response.LoginUser{
			ID:       user.ID,
			Email:    user.Email,
			Nickname: user.Nickname,
		},
		Token:        pair.AccessToken,
		RefreshToken: pair.RefreshToken,
	})
}

func (u *UserApi) Profile(ctx *gin.Context) {
	claims, err := utils.GetContextClaims(ctx)
	if err != nil || claims == nil {
		response.FailWithCode(ctx, response.UserNotExist)
		return
	}
	user, err := u.service.Profile(ctx, claims.User.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.SuccessWithData(ctx, &response.ProfileResponse{
		User: user,
		Roles: lo.Map(user.Roles, func(item *models.Role, _ int) *response.ProfileResponseRole {
			return &response.ProfileResponseRole{
				ID:   item.ID,
				Name: item.Name,
			}
		}),
		Departments: lo.Map(user.Departments, func(item *models.Department, _ int) *response.ProfileResponseDept {
			return &response.ProfileResponseDept{
				ID:   item.ID,
				Name: item.Name,
			}
		}),
	})
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
	response.SuccessWithPageData(ctx, total, lo.Map(list, func(item *models.User, _ int) *response.UserListRowResponse {
		return response.ToUserListRowResponse(item)
	}))
}
