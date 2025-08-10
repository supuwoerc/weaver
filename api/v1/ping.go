package v1

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/supuwoerc/weaver/pkg/response"
)

type PingService interface {
	LockPermissionField(ctx context.Context) error
}
type PingApi struct {
	*BasicApi
	service PingService
}

func NewPingApi(basic *BasicApi, service PingService) *PingApi {
	pinApi := &PingApi{
		BasicApi: basic,
		service:  service,
	}
	group := basic.route.Group("ping")
	{
		group.GET("", pinApi.Ping)
		group.GET("exception", pinApi.Exception)
		group.GET("check-permission", basic.auth.LoginRequired(), basic.auth.PermissionRequired(), pinApi.CheckPermission)
		group.GET("slow", pinApi.SlowResponse)
		group.GET("check-lock", pinApi.LockResponse)
		group.GET("logger-trace", pinApi.LoggerTrace)
	}
	return pinApi
}

// Ping
//
//	@Summary		健康检查
//	@Description	简单的健康检查接口，返回pong
//	@Tags			系统监控
//	@Accept			json
//	@Produce		json
//	@Success		10000	{object}	response.BasicResponse[string]	"健康检查成功，code=10000"
//	@Router			/ping [get]
func (p *PingApi) Ping(ctx *gin.Context) {
	response.SuccessWithData[string](ctx, "pong")
}

// Exception
//
//	@Summary		异常测试
//	@Description	故意触发除零异常进行测试
//	@Tags			系统监控
//	@Accept			json
//	@Produce		json
//	@Success		10000	{object}	response.BasicResponse[int]	"异常测试成功，code=10000"
//	@Router			/ping/exception [get]
func (p *PingApi) Exception(ctx *gin.Context) {
	num := 100 - (99 + 1)
	response.SuccessWithData[int](ctx, 1/num)
}

// CheckPermission
//
//	@Summary		权限检查测试
//	@Description	检查用户权限的测试接口
//	@Tags			系统监控
//	@Accept			json
//	@Produce		text/plain
//	@Security		BearerAuth
//	@Router			/ping/check-permission [get]
func (p *PingApi) CheckPermission(ctx *gin.Context) {
	ctx.String(http.StatusOK, "ok")
}

// SlowResponse
//
//	@Summary		慢响应测试
//	@Description	模拟慢响应的测试接口
//	@Tags			系统监控
//	@Accept			json
//	@Produce		text/plain
//	@Param			t	query	int	true	"睡眠秒数"
//	@Router			/ping/slow [get]
func (p *PingApi) SlowResponse(ctx *gin.Context) {
	value := ctx.Query("t")
	second, err := strconv.Atoi(value)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
	}
	time.Sleep(time.Duration(second) * time.Second)
	ctx.String(http.StatusOK, fmt.Sprintf("sleep %ds,PID %d", second, os.Getpid()))
}

// LockResponse
//
//	@Summary		锁测试
//	@Description	测试权限字段锁功能
//	@Tags			系统监控
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		10000	{object}	response.BasicResponse[any]	"锁测试成功，code=10000"
//	@Failure		10001	{object}	response.BasicResponse[any]	"锁测试失败，code=10001"
//	@Router			/ping/check-lock [get]
func (p *PingApi) LockResponse(ctx *gin.Context) {
	err := p.service.LockPermissionField(ctx)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}
	response.Success(ctx)
}

// LoggerTrace
//
//	@Summary		日志追踪测试
//	@Description	测试日志追踪功能
//	@Tags			系统监控
//	@Accept			json
//	@Produce		text/plain
//	@Router			/ping/logger-trace [get]
func (p *PingApi) LoggerTrace(ctx *gin.Context) {
	p.logger.WithContext(ctx).Infow("test message", "content", "hello trace!!!")
	ctx.String(http.StatusOK, "ok")
}
