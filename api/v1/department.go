package v1

import (
	"context"
	"gin-web/middleware"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"

	"github.com/gin-gonic/gin"
)

type DepartmentService interface {
	CreateDepartment(ctx context.Context, operator uint, params *request.CreateDepartmentRequest) error
	GetDepartmentTree(ctx context.Context, crew bool) ([]*response.DepartmentTreeResponse, error)
}

type DepartmentApi struct {
	service DepartmentService
}

func NewDepartmentApi(
	route *gin.RouterGroup,
	service DepartmentService,
	authMiddleware *middleware.AuthMiddleware,
) *DepartmentApi {
	departmentApi := &DepartmentApi{
		service: service,
	}
	// 挂载路由
	departmentAccessGroup := route.Group("department").Use(authMiddleware.LoginRequired())
	{
		departmentAccessGroup.POST("create", authMiddleware.PermissionRequired(), departmentApi.CreateDepartment)
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
	res, err := r.service.GetDepartmentTree(ctx, params.Crew)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.SuccessWithData(ctx, res)
}
