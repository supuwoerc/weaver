package v1

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/supuwoerc/weaver/middleware"
	"github.com/supuwoerc/weaver/pkg/response"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type PingService interface {
	LockPermissionField(ctx context.Context) error
}
type PingApi struct {
	service PingService
	logger  *zap.SugaredLogger
}

func NewPingApi(route *gin.RouterGroup, service PingService,
	authMiddleware *middleware.AuthMiddleware, logger *zap.SugaredLogger) *PingApi {
	pinApi := &PingApi{
		service: service,
		logger:  logger,
	}
	// 挂载路由
	group := route.Group("ping")
	{
		group.GET("", pinApi.Ping)
		group.GET("exception", pinApi.Exception)
		group.GET("check-permission", authMiddleware.LoginRequired(), authMiddleware.PermissionRequired(), pinApi.CheckPermission)
		group.GET("slow", pinApi.SlowResponse)
		group.GET("check-lock", pinApi.LockResponse)
		group.GET("logger-trace", pinApi.LoggerTrace)
	}
	return pinApi
}

func (p *PingApi) Ping(ctx *gin.Context) {
	response.SuccessWithData[string](ctx, "pong")
}

func (p *PingApi) Exception(ctx *gin.Context) {
	num := 100 - (99 + 1)
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
func (p *PingApi) LoggerTrace(ctx *gin.Context) {
	p.logger.Info("test trace", "这是测试内容", "99887766")
	ctx.String(http.StatusOK, "ok")
}
