package v1

import (
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/service"
	"github.com/gin-gonic/gin"
	"sync"
)

type RoleApi struct {
	*BasicApi
	service *service.RoleService
}

var (
	roleOnce sync.Once
	roleApi  *RoleApi
)

func NewRoleApi() *RoleApi {
	roleOnce.Do(func() {
		roleApi = &RoleApi{
			BasicApi: NewBasicApi(),
			service:  service.NewRoleService(),
		}
	})
	return roleApi
}

func (r *RoleApi) CreateRole(ctx *gin.Context) {
	var params request.CreateRoleRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	err := r.service.CreateRole(ctx, params.Name, params.Users, params.Permissions)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.Success(ctx)
}

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
