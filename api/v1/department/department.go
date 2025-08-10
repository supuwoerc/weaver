package department

import (
	"context"

	v1 "github.com/supuwoerc/weaver/api/v1"
	"github.com/supuwoerc/weaver/pkg/request"
	"github.com/supuwoerc/weaver/pkg/response"
	"github.com/supuwoerc/weaver/pkg/utils"

	"github.com/gin-gonic/gin"
)

type Service interface {
	CreateDepartment(ctx context.Context, operator uint, params *request.CreateDepartmentRequest) error
	GetDepartmentTree(ctx context.Context, withCrew bool) ([]*response.DepartmentTreeResponse, error)
}

type Api struct {
	*v1.BasicApi
	service Service
}

func NewDepartmentApi(basic *v1.BasicApi, service Service) *Api {
	departmentApi := &Api{
		BasicApi: basic,
		service:  service,
	}
	// 挂载路由
	departmentAccessGroup := basic.Route.Group("department").Use(basic.Auth.LoginRequired())
	{
		departmentAccessGroup.POST("create", basic.Auth.PermissionRequired(), departmentApi.CreateDepartment)
		departmentAccessGroup.GET("tree", departmentApi.GetDepartmentTree)
	}
	return departmentApi
}

// CreateDepartment
//
//	@Summary		创建部门
//	@Description	创建新的部门
//	@Tags			部门管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		request.CreateDepartmentRequest	true	"创建部门请求参数"
//	@Success		10000	{object}	response.BasicResponse[any]		"创建成功，code=10000"
//	@Failure		10002	{object}	response.BasicResponse[any]		"参数验证失败，code=10002"
//	@Failure		10008	{object}	response.BasicResponse[any]		"认证失败，code=10008"
//	@Failure		10001	{object}	response.BasicResponse[any]		"业务逻辑失败，code=10001"
//	@Router			/department/create [post]
func (r *Api) CreateDepartment(ctx *gin.Context) {
	var params request.CreateDepartmentRequest
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
	err = r.service.CreateDepartment(ctx, claims.User.ID, &params)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.Success(ctx)
}

// GetDepartmentTree
//
//	@Summary		获取部门树
//	@Description	获取部门树形结构，可选择是否包含人员信息
//	@Tags			部门管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			withCrew	query		bool														false	"是否包含人员信息"	default(false)
//	@Success		10000		{object}	response.BasicResponse[[]response.DepartmentTreeResponse]	"获取成功，code=10000"
//	@Failure		10002		{object}	response.BasicResponse[any]									"参数验证失败，code=10002"
//	@Failure		10001		{object}	response.BasicResponse[any]									"服务器内部错误，code=10001"
//	@Router			/department/tree [get]
func (r *Api) GetDepartmentTree(ctx *gin.Context) {
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
