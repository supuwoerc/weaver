package v1

import (
	"gin-web/models"
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
	departments, err := r.service.GetDepartmentTree(ctx, params.Crew)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	var res []*response.DepartmentTreeResponse
	nodeMap := make(map[uint]*response.DepartmentTreeResponse)
	deptMap := make(map[uint]*models.Department)
	for _, dept := range departments {
		deptMap[dept.ID] = dept
	}
	for _, dept := range departments {
		holder, exist := nodeMap[dept.ID]
		var children = make([]*response.DepartmentTreeResponse, 0)
		if exist {
			children = holder.Children
		}
		node, parseErr := response.ToDepartmentTreeResponse(dept, deptMap)
		if parseErr != nil {
			response.FailWithError(ctx, parseErr)
			return
		}
		node.Children = children
		nodeMap[node.ID] = node
		if dept.ParentId == nil {
			res = append(res, node)
		} else {
			_, exist = nodeMap[*dept.ParentId]
			if !exist {
				nodeMap[*dept.ParentId], parseErr = response.ToDepartmentTreeResponse(&models.Department{}, deptMap)
				if parseErr != nil {
					response.FailWithError(ctx, parseErr)
					return
				}
			}
			nodeMap[*dept.ParentId].Children = append(nodeMap[*dept.ParentId].Children, node)
		}
	}
	response.SuccessWithData(ctx, res)
}
