package v1

import (
	"context"
	"gin-web/pkg/redis"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"gin-web/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
)

type DepartmentService interface {
	CreateDepartment(ctx context.Context, operator uint, name string, parentId *uint, leaderIds, userIds []uint) error
	GetDepartmentTree(ctx context.Context, crew bool) ([]*response.DepartmentTreeResponse, error)
}

type DepartmentApi struct {
	*BasicApi
	service DepartmentService
}

var (
	departmentOnce sync.Once
	departmentApi  *DepartmentApi
)

func NewDepartmentApi(route *gin.RouterGroup, logger *zap.SugaredLogger, r *redis.CommonRedisClient, db *gorm.DB,
	locksmith *utils.RedisLocksmith, v *viper.Viper) *DepartmentApi {
	departmentOnce.Do(func() {
		departmentApi = &DepartmentApi{
			BasicApi: NewBasicApi(logger, v),
			service:  service.NewDepartmentService(logger, db, r, locksmith, v),
		}
		// 挂载路由
		departmentAccessGroup := route.Group("department")
		{
			departmentAccessGroup.POST("create", departmentApi.CreateDepartment)
			departmentAccessGroup.GET("tree", departmentApi.GetDepartmentTree)
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
