package v1

import (
	"context"

	"github.com/supuwoerc/weaver/pkg/request"
	"github.com/supuwoerc/weaver/pkg/response"
	"github.com/supuwoerc/weaver/pkg/utils"

	"github.com/gin-gonic/gin"
)

type PermissionService interface {
	CreatePermission(ctx context.Context, operator uint, params *request.CreatePermissionRequest) error
	GetPermissionList(ctx context.Context, keyword string, limit, offset int) ([]*response.PermissionListRowResponse, int64, error)
	GetPermissionDetail(ctx context.Context, id uint) (*response.PermissionDetailResponse, error)
	UpdatePermission(ctx context.Context, operator uint, params *request.UpdatePermissionRequest) error
	DeletePermission(ctx context.Context, id, operator uint) error
}

type PermissionApi struct {
	*BasicApi
	service PermissionService
}

func NewPermissionApi(basic *BasicApi, service PermissionService) *PermissionApi {
	permissionApi := &PermissionApi{
		BasicApi: basic,
		service:  service,
	}
	// 挂载路由
	permissionAccessGroup := basic.route.Group("permission").Use(
		basic.auth.LoginRequired(),
		basic.auth.PermissionRequired(),
	)
	{
		permissionAccessGroup.POST("create", permissionApi.CreatePermission)
		permissionAccessGroup.GET("list", permissionApi.GetPermissionList)
		permissionAccessGroup.GET("detail", permissionApi.GetPermissionDetail)
		permissionAccessGroup.POST("update", permissionApi.UpdatePermission)
		permissionAccessGroup.POST("delete", permissionApi.DeletePermission)
	}
	return permissionApi
}

func (r *PermissionApi) CreatePermission(ctx *gin.Context) {
	var params request.CreatePermissionRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
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

func (r *PermissionApi) GetPermissionList(ctx *gin.Context) {
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

func (r *PermissionApi) GetPermissionDetail(ctx *gin.Context) {
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

func (r *PermissionApi) UpdatePermission(ctx *gin.Context) {
	var params request.UpdatePermissionRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
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

func (r *PermissionApi) DeletePermission(ctx *gin.Context) {
	var params request.DeletePermissionRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
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
