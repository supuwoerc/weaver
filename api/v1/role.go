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

type RoleService interface {
	CreateRole(ctx context.Context, operator uint, name string, userIds, permissionIds []uint) error
	GetRoleList(ctx context.Context, keyword string, limit, offset int) ([]*response.RoleListRowResponse, int64, error)
	GetRoleDetail(ctx context.Context, id uint) (*response.RoleDetailResponse, error)
	UpdateRole(ctx context.Context, operator uint, id uint, name string, userIds, permissionIds []uint) error
	DeleteRole(ctx context.Context, id, operator uint) error
}

type RoleApi struct {
	*BasicApi
	service RoleService
}

var (
	roleOnce sync.Once
	roleApi  *RoleApi
)

func NewRoleApi(route *gin.RouterGroup, logger *zap.SugaredLogger, r *redis.CommonRedisClient, db *gorm.DB,
	locksmith *utils.RedisLocksmith, v *viper.Viper) *RoleApi {
	roleOnce.Do(func() {
		roleApi = &RoleApi{
			BasicApi: NewBasicApi(logger, v),
			service:  service.NewRoleService(logger, r, db, locksmith, v),
		}
		// 挂载路由
		roleAccessGroup := route.Group("role")
		{
			roleAccessGroup.POST("create", roleApi.CreateRole)
			roleAccessGroup.GET("list", roleApi.GetRoleList)
			roleAccessGroup.GET("detail", roleApi.GetRoleDetail)
			roleAccessGroup.POST("update", roleApi.UpdateRole)
			roleAccessGroup.POST("delete", roleApi.DeleteRole)
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
