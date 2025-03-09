package v1

import (
	"context"
	"fmt"
	"gin-web/middleware"
	"gin-web/pkg/redis"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"gin-web/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type PingService interface {
	LockPermissionField(ctx context.Context) error
}
type PingApi struct {
	*BasicApi
	service PingService
}

var (
	pingOnce sync.Once
	pinApi   *PingApi
)

func NewPingApi(route *gin.RouterGroup, logger *zap.SugaredLogger, r *redis.CommonRedisClient, db *gorm.DB, redisClient *redis.CommonRedisClient,
	locksmith *utils.RedisLocksmith, v *viper.Viper) *PingApi {
	pingOnce.Do(func() {
		pinApi = &PingApi{
			BasicApi: NewBasicApi(logger, v),
			service:  service.NewPingService(logger, db, r, locksmith, v),
		}
		// 挂载路由
		group := route.Group("ping")
		{
			group.GET("", pinApi.Ping)
			group.GET("exception", pinApi.Exception)
			group.GET("check-permission", middleware.NewAuthMiddleware(db, redisClient, v).PermissionRequired(), pinApi.CheckPermission)
			group.GET("slow", pinApi.SlowResponse)
			group.GET("check-lock", pinApi.LockResponse)
		}
	})
	return pinApi
}

func (p *PingApi) Ping(ctx *gin.Context) {
	response.SuccessWithData[string](ctx, "pong")
}

func (p *PingApi) Exception(ctx *gin.Context) {
	num := 100 - 100
	response.SuccessWithData[int](ctx, 1/num)
}

func (p *PingApi) CheckPermission(ctx *gin.Context) {
	ctx.String(http.StatusOK, "ok")
}
func (p *PingApi) SlowResponse(ctx *gin.Context) {
	value := ctx.Query("t")
	second, err := strconv.Atoi(value)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
	}
	time.Sleep(time.Duration(second) * time.Second)
	ctx.String(http.StatusOK, fmt.Sprintf("sleep %ds,PID %d", second, os.Getpid()))
}
func (p *PingApi) LockResponse(ctx *gin.Context) {
	err := p.service.LockPermissionField(ctx)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.Success(ctx)
}
