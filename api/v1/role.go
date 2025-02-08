package v1

import (
	"gin-web/models"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"gin-web/service"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
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
	claims, err := utils.GetContextClaims(ctx)
	if err != nil || claims == nil {
		response.FailWithCode(ctx, response.AuthErr)
		return
	}
	err = r.service.CreateRole(ctx, claims.User.ID, params.Name, params.Users, params.Permissions)
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
	response.SuccessWithPageData(ctx, total, lo.Map(list, func(item *models.Role, _ int) *response.RoleListRowResponse {
		return response.ToRoleListRowResponse(item)
	}))
}

func (r *RoleApi) GetRoleDetail(ctx *gin.Context) {
	var params request.GetRoleDetailRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	role, err := r.service.GetRoleDetail(ctx, params.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.SuccessWithData(ctx, response.ToRoleDetailResponse(role))
}

func (r *RoleApi) UpdateRole(ctx *gin.Context) {
	var params request.UpdateRoleRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	claims, err := utils.GetContextClaims(ctx)
	if err != nil || claims == nil {
		response.FailWithCode(ctx, response.AuthErr)
		return
	}
	err = r.service.UpdateRole(ctx, claims.User.ID, params.ID, params.Name, params.Users, params.Permissions)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.Success(ctx)
}

func (r *RoleApi) DeleteRole(ctx *gin.Context) {
	var params request.DeleteRoleRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
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
