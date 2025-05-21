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
	p.logger.WithContext(ctx).Infow("test message", "content", "hello trace!!!")
	ctx.String(http.StatusOK, "ok")
}
