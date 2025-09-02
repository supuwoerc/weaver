package user

import (
	"context"
	"net/http"

	v1 "github.com/supuwoerc/weaver/api/v1"
	"github.com/supuwoerc/weaver/models"
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/request"
	"github.com/supuwoerc/weaver/pkg/response"
	"github.com/supuwoerc/weaver/pkg/utils"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
)

type Service interface {
	SignUp(ctx context.Context, id string, code string, user *models.User) error
	Login(ctx context.Context, email string, password string) (*response.LoginResponse, error)
	Logout(ctx context.Context, email string) error
	Profile(ctx context.Context, uid uint) (*response.ProfileResponse, error)
	GetUserList(ctx context.Context, keyword string, limit, offset int) ([]*response.UserListRowResponse, int64, error)
	ActiveAccount(ctx context.Context, uid uint, activeCode string) error
}

type Api struct {
	*v1.BasicApi
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
	phoneRegexExp    *regexp.Regexp
	service          Service
}

func NewUserApi(basic *v1.BasicApi, service Service) *Api {
	userApi := &Api{
		BasicApi:         basic,
		emailRegexExp:    regexp.MustCompile(constant.EmailRegexPattern, regexp.None),
		passwordRegexExp: regexp.MustCompile(constant.PasswdRegexPattern, regexp.None),
		phoneRegexExp:    regexp.MustCompile(constant.PhoneRegexPattern, regexp.None),
		service:          service,
	}
	// 挂载路由
	userPublicGroup := basic.Route.Group("public/user")
	{
		userPublicGroup.POST("signup", userApi.SignUp)
		userPublicGroup.POST("login", userApi.Login)
		userPublicGroup.GET("active", userApi.Active)
		userPublicGroup.GET("active-success", userApi.ActiveSuccess)
		userPublicGroup.GET("active-failure", userApi.ActiveFailure)
	}
	userAccessGroup := basic.Route.Group("user").Use(basic.Auth.LoginRequired())
	{
		userAccessGroup.POST("refresh-token")
		userAccessGroup.GET("profile", userApi.Profile)
		userAccessGroup.POST("logout", userApi.Logout)
		userAccessGroup.GET("list", basic.Auth.PermissionRequired(), userApi.GetUserList)
	}
	return userApi
}

// SignUp
//
//	@Summary		用户注册
//	@Description	用户通过邮箱和密码进行注册
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.SignUpRequest		true	"注册请求参数"
//	@Success		10000	{object}	response.BasicResponse[any]	"注册成功，code=10000"
//	@Failure		10002	{object}	response.BasicResponse[any]	"参数验证失败，code=10002"
//	@Failure		20004	{object}	response.BasicResponse[any]	"邮箱格式错误，code=20004"
//	@Failure		20003	{object}	response.BasicResponse[any]	"密码格式错误，code=20003"
//	@Failure		10001	{object}	response.BasicResponse[any]	"业务逻辑失败，code=10001"
//	@Router			/public/user/signup [post]
func (r *Api) SignUp(ctx *gin.Context) {
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

// Login
//
//	@Summary		用户登录
//	@Description	用户通过邮箱和密码进行登录
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.LoginRequest							true	"登录请求参数"
//	@Success		10000	{object}	response.BasicResponse[response.LoginResponse]	"登录成功，code=10000"
//	@Failure		10002	{object}	response.BasicResponse[any]						"参数验证失败，code=10002"
//	@Failure		20001	{object}	response.BasicResponse[any]						"登录失败，code=20001"
//	@Router			/public/user/login [post]
func (r *Api) Login(ctx *gin.Context) {
	var params request.LoginRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	res, err := r.service.Login(ctx, params.Email, params.Password)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.SuccessWithData(ctx, res)
}

// Logout
//
//	@Summary		退出登录
//	@Description	用户登出(退出登录)
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Success		10000	{object}	response.BasicResponse[response.LoginResponse]	"登出成功，code=10000"
//	@Failure		20001	{object}	response.BasicResponse[any]						"登出失败，code=20001"
//	@Router			/user/logout [post]
func (r *Api) Logout(ctx *gin.Context) {
	claims, err := utils.GetContextClaims(ctx)
	if err != nil || claims == nil {
		response.FailWithCode(ctx, response.UserNotExist)
		return
	}
	err = r.service.Logout(ctx, claims.User.Email)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.Success(ctx)
}

// Profile
//
//	@Summary		获取用户资料
//	@Description	获取当前登录用户的详细资料信息
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		10000	{object}	response.BasicResponse[response.ProfileResponse]	"获取成功，code=10000"
//	@Failure		20005	{object}	response.BasicResponse[any]							"用户不存在，code=20005"
//	@Failure		10001	{object}	response.BasicResponse[any]							"服务器内部错误，code=10001"
//	@Router			/user/profile [get]
func (r *Api) Profile(ctx *gin.Context) {
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

// GetUserList
//
//	@Summary		获取用户列表
//	@Description	分页获取用户列表，支持关键词搜索
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			keyword	query		string																	false	"搜索关键词"
//	@Param			limit	query		int																		false	"每页数量"	default(10)
//	@Param			offset	query		int																		false	"偏移量"	default(0)
//	@Success		10000	{object}	response.BasicResponse[response.DataList[response.UserListRowResponse]]	"获取成功，code=10000"
//	@Failure		10002	{object}	response.BasicResponse[any]												"参数验证失败，code=10002"
//	@Failure		10001	{object}	response.BasicResponse[any]												"服务器内部错误，code=10001"
//	@Router			/user/list [get]
func (r *Api) GetUserList(ctx *gin.Context) {
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

func (r *Api) Active(ctx *gin.Context) {
	var params request.ActiveAccountRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.Redirect(http.StatusMovedPermanently, "/view/v1/public/user/active-failure")
		return
	}
	err := r.service.ActiveAccount(ctx, params.ID, params.ActiveCode)
	if err != nil {
		ctx.Redirect(http.StatusMovedPermanently, "/view/v1/public/user/active-failure")
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/view/v1/public/user/active-success")
	}
}
func (r *Api) ActiveSuccess(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "active-success.html", nil)
}

func (r *Api) ActiveFailure(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "active-failure.html", nil)
}
