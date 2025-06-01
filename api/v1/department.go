package v1

import (
	"context"

	"github.com/supuwoerc/weaver/pkg/request"
	"github.com/supuwoerc/weaver/pkg/response"
	"github.com/supuwoerc/weaver/pkg/utils"

	"github.com/gin-gonic/gin"
)

type DepartmentService interface {
	CreateDepartment(ctx context.Context, operator uint, params *request.CreateDepartmentRequest) error
	GetDepartmentTree(ctx context.Context, withCrew bool) ([]*response.DepartmentTreeResponse, error)
}

type DepartmentApi struct {
	*BasicApi
	service DepartmentService
}

func NewDepartmentApi(basic *BasicApi, service DepartmentService) *DepartmentApi {
	departmentApi := &DepartmentApi{
		BasicApi: basic,
		service:  service,
	}
	// 挂载路由
	departmentAccessGroup := basic.route.Group("department").Use(basic.auth.LoginRequired())
	{
		departmentAccessGroup.POST("create", basic.auth.PermissionRequired(), departmentApi.CreateDepartment)
		departmentAccessGroup.GET("tree", departmentApi.GetDepartmentTree)
	}
	return departmentApi
}

func (r *DepartmentApi) CreateDepartment(ctx *gin.Context) {
	var params request.CreateDepartmentRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	claims, err := utils.GetContextClaims(ctx)
	if err != nil || claims == nil {
		response.FailWithCode(ctx, response.AuthErr)
		return
	}
	err = r.service.CreateDepartment(ctx, claims.User.ID, &params)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.Success(ctx)
}

func (r *DepartmentApi) GetDepartmentTree(ctx *gin.Context) {
	var params request.GetDepartmentTreeRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.ParamsValidateFail(ctx, err)
		return
	}
	res, err := r.service.GetDepartmentTree(ctx, params.WithCrew)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.SuccessWithData(ctx, res)
}
