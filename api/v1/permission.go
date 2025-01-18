package v1

import (
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/service"
	"github.com/gin-gonic/gin"
	"sync"
)

type PermissionApi struct {
	*BasicApi
	service *service.PermissionService
}

var (
	permissionOnce sync.Once
	permissionApi  *PermissionApi
)

func NewPermissionApi() *PermissionApi {
	permissionOnce.Do(func() {
		permissionApi = &PermissionApi{
			BasicApi: NewBasicApi(),
			service:  service.NewPermissionService(),
		}
	})
	return permissionApi
}

func (r *PermissionApi) CreatePermission(ctx *gin.Context) {
	var params request.CreatePermissionRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	err := r.service.CreatePermission(ctx, params.Name, params.Resource, params.Roles)
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
