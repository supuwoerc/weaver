package permission

import (
	"context"

	v1 "github.com/supuwoerc/weaver/api/v1"
	"github.com/supuwoerc/weaver/pkg/request"
	"github.com/supuwoerc/weaver/pkg/response"
	"github.com/supuwoerc/weaver/pkg/utils"

	"github.com/gin-gonic/gin"
)

type Service interface {
	CreatePermission(ctx context.Context, operator uint, params *request.CreatePermissionRequest) error
	GetPermissionList(ctx context.Context, keyword string, limit, offset int) ([]*response.PermissionListRowResponse, int64, error)
	GetPermissionDetail(ctx context.Context, id uint) (*response.PermissionDetailResponse, error)
	UpdatePermission(ctx context.Context, operator uint, params *request.UpdatePermissionRequest) error
	DeletePermission(ctx context.Context, id, operator uint) error
	GetUserViewRouteAndMenuPermissions(ctx context.Context, uid uint) (response.FrontEndPermissions, error)
}

type Api struct {
	*v1.BasicApi
	service Service
}

func NewPermissionApi(basic *v1.BasicApi, service Service) *Api {
	permissionApi := &Api{
		BasicApi: basic,
		service:  service,
	}
	// 挂载路由
	permissionGroup := basic.Route.Group("permission").Use(basic.Auth.LoginRequired())
	{
		permissionGroup.GET("user-route-menu-permissions", permissionApi.GetUserViewRouteAndMenuPermissions)
	}
	permissionAccessGroup := permissionGroup.Use(basic.Auth.PermissionRequired())
	{
		permissionAccessGroup.POST("create", permissionApi.CreatePermission)
		permissionAccessGroup.GET("list", permissionApi.GetPermissionList)
		permissionAccessGroup.GET("detail", permissionApi.GetPermissionDetail)
		permissionAccessGroup.POST("update", permissionApi.UpdatePermission)
		permissionAccessGroup.POST("delete", permissionApi.DeletePermission)
	}
	return permissionApi
}

// CreatePermission
//
//	@Summary		创建权限
//	@Description	创建新的权限
//	@Tags			权限管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		request.CreatePermissionRequest	true	"创建权限请求参数"
//	@Success		10000	{object}	response.BasicResponse[any]		"创建成功，code=10000"
//	@Failure		10002	{object}	response.BasicResponse[any]		"参数验证失败，code=10002"
//	@Failure		10008	{object}	response.BasicResponse[any]		"认证失败，code=10008"
//	@Failure		10001	{object}	response.BasicResponse[any]		"业务逻辑失败，code=10001"
//	@Router			/permission/create [post]
func (r *Api) CreatePermission(ctx *gin.Context) {
	var params request.CreatePermissionRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	// TODO:集中到鉴权中间件
	claims, err := utils.GetContextClaims(ctx)
	if err != nil || claims == nil {
		response.FailWithCode(ctx, response.AuthErr)
		return
	}
	err = r.service.CreatePermission(ctx, claims.User.ID, &params)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.Success(ctx)
}

// GetPermissionList
//
//	@Summary		获取权限列表
//	@Description	分页获取权限列表，支持关键词搜索
//	@Tags			权限管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			keyword	query		string																			false	"搜索关键词"
//	@Param			limit	query		int																				false	"每页数量"	default(10)
//	@Param			offset	query		int																				false	"偏移量"	default(0)
//	@Success		10000	{object}	response.BasicResponse[response.DataList[response.PermissionListRowResponse]]	"获取成功，code=10000"
//	@Failure		10002	{object}	response.BasicResponse[any]														"参数验证失败，code=10002"
//	@Failure		10001	{object}	response.BasicResponse[any]														"服务器内部错误，code=10001"
//	@Router			/permission/list [get]
func (r *Api) GetPermissionList(ctx *gin.Context) {
	var params request.GetPermissionListRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	list, total, err := r.service.GetPermissionList(ctx, params.Keyword, params.Limit, params.Offset)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.SuccessWithPageData(ctx, total, list)
}

// GetPermissionDetail
//
//	@Summary		获取权限详情
//	@Description	根据权限ID获取权限详细信息
//	@Tags			权限管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		query		uint														true	"权限ID"
//	@Success		10000	{object}	response.BasicResponse[response.PermissionDetailResponse]	"获取成功，code=10000"
//	@Failure		10002	{object}	response.BasicResponse[any]									"参数验证失败，code=10002"
//	@Failure		10001	{object}	response.BasicResponse[any]									"服务器内部错误，code=10001"
//	@Router			/permission/detail [get]
func (r *Api) GetPermissionDetail(ctx *gin.Context) {
	var params request.GetPermissionDetailRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	detail, err := r.service.GetPermissionDetail(ctx, params.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.SuccessWithData(ctx, detail)
}

// UpdatePermission
//
//	@Summary		更新权限
//	@Description	更新权限信息
//	@Tags			权限管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		request.UpdatePermissionRequest	true	"更新权限请求参数"
//	@Success		10000	{object}	response.BasicResponse[any]		"更新成功，code=10000"
//	@Failure		10002	{object}	response.BasicResponse[any]		"参数验证失败，code=10002"
//	@Failure		10008	{object}	response.BasicResponse[any]		"认证失败，code=10008"
//	@Failure		10001	{object}	response.BasicResponse[any]		"业务逻辑失败，code=10001"
//	@Router			/permission/update [post]
func (r *Api) UpdatePermission(ctx *gin.Context) {
	var params request.UpdatePermissionRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	// TODO:集中到鉴权中间件
	claims, err := utils.GetContextClaims(ctx)
	if err != nil || claims == nil {
		response.FailWithCode(ctx, response.AuthErr)
		return
	}
	err = r.service.UpdatePermission(ctx, claims.User.ID, &params)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.Success(ctx)
}

// DeletePermission
//
//	@Summary		删除权限
//	@Description	根据权限ID删除权限
//	@Tags			权限管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		request.DeletePermissionRequest	true	"删除权限请求参数"
//	@Success		10000	{object}	response.BasicResponse[any]		"删除成功，code=10000"
//	@Failure		10002	{object}	response.BasicResponse[any]		"参数验证失败，code=10002"
//	@Failure		10008	{object}	response.BasicResponse[any]		"认证失败，code=10008"
//	@Failure		10001	{object}	response.BasicResponse[any]		"业务逻辑失败，code=10001"
//	@Router			/permission/delete [post]
func (r *Api) DeletePermission(ctx *gin.Context) {
	var params request.DeletePermissionRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	// TODO:集中到鉴权中间件
	claims, err := utils.GetContextClaims(ctx)
	if err != nil || claims == nil {
		response.FailWithCode(ctx, response.AuthErr)
		return
	}
	err = r.service.DeletePermission(ctx, params.ID, claims.User.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.Success(ctx)
}

// GetUserViewRouteAndMenuPermissions
//
//	@Summary		获取账户可访问的前端权限(菜单权限 & 路由权限)
//	@Description	获取账户可访问的前端权限(菜单权限 & 路由权限)
//	@Tags			权限管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		10000	{object}	response.BasicResponse[response.FrontEndPermissions]	"获取成功，code=10000"
//	@Failure		10002	{object}	response.BasicResponse[any]								"参数验证失败，code=10002"
//	@Failure		10001	{object}	response.BasicResponse[any]								"服务器内部错误，code=10001"
//	@Router			/permission/user-route-menu-permissions [get]
func (r *Api) GetUserViewRouteAndMenuPermissions(ctx *gin.Context) {
	claims, err := utils.GetContextClaims(ctx)
	if err != nil || claims == nil {
		response.FailWithCode(ctx, response.AuthErr)
		return
	}
	list, err := r.service.GetUserViewRouteAndMenuPermissions(ctx, claims.User.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.SuccessWithData(ctx, list)
}
