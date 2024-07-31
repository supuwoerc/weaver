package v1

import (
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/service"
	"github.com/gin-gonic/gin"
)

type RoleApi struct {
	*BasicApi
	service func(ctx *gin.Context) *service.RoleService
}

func NewRoleApi() RoleApi {
	return RoleApi{
		BasicApi: NewBasicApi(),
		service: func(ctx *gin.Context) *service.RoleService {
			return service.NewRoleService(ctx)
		},
	}
}

// TODO:添加swagger文档注释
func (r RoleApi) CreateRole(ctx *gin.Context) {
	var params request.CreateRoleRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	roleService := service.NewRoleService(ctx)
	err := roleService.CreateRole(ctx, params.Name)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.Success(ctx)
}
