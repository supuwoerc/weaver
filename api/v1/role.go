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
	err := r.service.CreateRole(ctx, params.Name)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.Success(ctx)
}

// TODO:添加swagger文档注释
func (r *RoleApi) GetRoleList(ctx *gin.Context) {
	var params request.GetRoleListRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	list, total, err := r.service.GetRoleList(ctx, params.Name, params.Limit, params.Offset)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.SuccessWithPageData(ctx, total, list)
}
