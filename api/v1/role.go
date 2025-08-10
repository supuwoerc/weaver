package v1

import (
	"context"

	"github.com/supuwoerc/weaver/pkg/request"
	"github.com/supuwoerc/weaver/pkg/response"
	"github.com/supuwoerc/weaver/pkg/utils"

	"github.com/gin-gonic/gin"
)

type RoleService interface {
	CreateRole(ctx context.Context, operator uint, params *request.CreateRoleRequest) error
	GetRoleList(ctx context.Context, keyword string, limit, offset int) ([]*response.RoleListRowResponse, int64, error)
	GetRoleDetail(ctx context.Context, id uint) (*response.RoleDetailResponse, error)
	UpdateRole(ctx context.Context, operator uint, params *request.UpdateRoleRequest) error
	DeleteRole(ctx context.Context, id, operator uint) error
}

type RoleApi struct {
	*BasicApi
	service RoleService
}

func NewRoleApi(basic *BasicApi, service RoleService) *RoleApi {
	roleApi := &RoleApi{
		BasicApi: basic,
		service:  service,
	}
	roleAccessGroup := basic.route.Group("role").Use(
		basic.auth.LoginRequired(),
		basic.auth.PermissionRequired(),
	)
	{
		roleAccessGroup.POST("create", roleApi.CreateRole)
		roleAccessGroup.GET("list", roleApi.GetRoleList)
		roleAccessGroup.GET("detail", roleApi.GetRoleDetail)
		roleAccessGroup.POST("update", roleApi.UpdateRole)
		roleAccessGroup.POST("delete", roleApi.DeleteRole)
	}
	return roleApi
}

// CreateRole
//
//	@Summary		创建角色
//	@Description	创建新的角色
//	@Tags			角色管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		request.CreateRoleRequest	true	"创建角色请求参数"
//	@Success		10000	{object}	response.BasicResponse[any]	"创建成功，code=10000"
//	@Failure		10002	{object}	response.BasicResponse[any]	"参数验证失败，code=10002"
//	@Failure		10008	{object}	response.BasicResponse[any]	"认证失败，code=10008"
//	@Failure		10001	{object}	response.BasicResponse[any]	"业务逻辑失败，code=10001"
//	@Router			/role/create [post]
func (r *RoleApi) CreateRole(ctx *gin.Context) {
	var params request.CreateRoleRequest
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
	err = r.service.CreateRole(ctx, claims.User.ID, &params)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.Success(ctx)
}

// GetRoleList
//
//	@Summary		获取角色列表
//	@Description	分页获取角色列表，支持关键词搜索
//	@Tags			角色管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			keyword	query		string																	false	"搜索关键词"
//	@Param			limit	query		int																		false	"每页数量"	default(10)
//	@Param			offset	query		int																		false	"偏移量"	default(0)
//	@Success		10000	{object}	response.BasicResponse[response.DataList[response.RoleListRowResponse]]	"获取成功，code=10000"
//	@Failure		10002	{object}	response.BasicResponse[any]												"参数验证失败，code=10002"
//	@Failure		10001	{object}	response.BasicResponse[any]												"服务器内部错误，code=10001"
//	@Router			/role/list [get]
func (r *RoleApi) GetRoleList(ctx *gin.Context) {
	var params request.GetRoleListRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	list, total, err := r.service.GetRoleList(ctx, params.Keyword, params.Limit, params.Offset)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.SuccessWithPageData(ctx, total, list)
}

// GetRoleDetail
//
//	@Summary		获取角色详情
//	@Description	根据角色ID获取角色详细信息
//	@Tags			角色管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		query		uint												true	"角色ID"
//	@Success		10000	{object}	response.BasicResponse[response.RoleDetailResponse]	"获取成功，code=10000"
//	@Failure		10002	{object}	response.BasicResponse[any]							"参数验证失败，code=10002"
//	@Failure		10001	{object}	response.BasicResponse[any]							"服务器内部错误，code=10001"
//	@Router			/role/detail [get]
func (r *RoleApi) GetRoleDetail(ctx *gin.Context) {
	var params request.GetRoleDetailRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	detail, err := r.service.GetRoleDetail(ctx, params.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.SuccessWithData(ctx, detail)
}

// UpdateRole
//
//	@Summary		更新角色
//	@Description	更新角色信息
//	@Tags			角色管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		request.UpdateRoleRequest	true	"更新角色请求参数"
//	@Success		10000	{object}	response.BasicResponse[any]	"更新成功，code=10000"
//	@Failure		10002	{object}	response.BasicResponse[any]	"参数验证失败，code=10002"
//	@Failure		10008	{object}	response.BasicResponse[any]	"认证失败，code=10008"
//	@Failure		10001	{object}	response.BasicResponse[any]	"业务逻辑失败，code=10001"
//	@Router			/role/update [post]
func (r *RoleApi) UpdateRole(ctx *gin.Context) {
	var params request.UpdateRoleRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	// TODO:集中到鉴权中间件中
	claims, err := utils.GetContextClaims(ctx)
	if err != nil || claims == nil {
		response.FailWithCode(ctx, response.AuthErr)
		return
	}
	err = r.service.UpdateRole(ctx, claims.User.ID, &params)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.Success(ctx)
}

// DeleteRole
//
//	@Summary		删除角色
//	@Description	根据角色ID删除角色
//	@Tags			角色管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		request.DeleteRoleRequest	true	"删除角色请求参数"
//	@Success		10000	{object}	response.BasicResponse[any]	"删除成功，code=10000"
//	@Failure		10002	{object}	response.BasicResponse[any]	"参数验证失败，code=10002"
//	@Failure		10008	{object}	response.BasicResponse[any]	"认证失败，code=10008"
//	@Failure		10001	{object}	response.BasicResponse[any]	"业务逻辑失败，code=10001"
//	@Router			/role/delete [post]
func (r *RoleApi) DeleteRole(ctx *gin.Context) {
	var params request.DeleteRoleRequest
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
	err = r.service.DeleteRole(ctx, params.ID, claims.User.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.Success(ctx)
}
