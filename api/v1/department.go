package v1

import (
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"gin-web/service"
	"github.com/gin-gonic/gin"
	"sync"
)

type DepartmentApi struct {
	*BasicApi
	service *service.DepartmentService
}

var (
	departmentOnce sync.Once
	departmentApi  *DepartmentApi
)

func NewDepartmentApi() *DepartmentApi {
	departmentOnce.Do(func() {
		departmentApi = &DepartmentApi{
			BasicApi: NewBasicApi(),
			service:  service.NewDepartmentService(),
		}
	})
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
	err = r.service.CreateDepartment(ctx, claims.User.ID, params.Name, params.ParentId, params.Leaders, params.Users)
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
