package ping

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	v1 "github.com/supuwoerc/weaver/api/v1"
	"github.com/supuwoerc/weaver/pkg/response"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Service interface {
	LockPermissionField(ctx context.Context) error
}
type Api struct {
	*v1.BasicApi
	service Service
}

func NewPingApi(basic *v1.BasicApi, service Service) *Api {
	pinApi := &Api{
		BasicApi: basic,
		service:  service,
	}
	group := basic.Route.Group("check")
	{
		group.GET("exception", pinApi.Exception)
		group.GET("check-permission", basic.Auth.LoginRequired(), basic.Auth.PermissionRequired(), pinApi.CheckPermission)
		group.GET("slow", pinApi.SlowResponse)
		group.GET("check-lock", pinApi.LockResponse)
		group.GET("logger-trace", pinApi.LoggerTrace)
		group.GET("span-trace", pinApi.SpanTrace)
	}
	return pinApi
}

// Exception
//
//	@Summary		异常测试
//	@Description	故意触发除零异常进行测试
//	@Tags			系统监控
//	@Accept			json
//	@Produce		json
//	@Success		10000	{object}	response.BasicResponse[int]	"异常测试成功，code=10000"
//	@Router			/check/exception [get]
func (p *Api) Exception(ctx *gin.Context) {
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
//	@Router			/check/check-permission [get]
func (p *Api) CheckPermission(ctx *gin.Context) {
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
//	@Router			/check/slow [get]
func (p *Api) SlowResponse(ctx *gin.Context) {
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
//	@Router			/check/check-lock [get]
func (p *Api) LockResponse(ctx *gin.Context) {
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
//	@Router			/check/logger-trace [get]
func (p *Api) LoggerTrace(ctx *gin.Context) {
	p.Logger.WithContext(ctx).Infow("test message", "content", "hello trace!!!")
	ctx.String(http.StatusOK, "ok")
}

// SpanTrace
//
//	@Summary		链路追踪
//	@Description	简单的链路追踪接口，返回trace
//	@Tags			系统监控
//	@Accept			json
//	@Produce		text/plain
//	@Success		10000	{object}	response.BasicResponse[string]	"链路追踪，code=10000"
//	@Router			/check/span-trace [get]
func (p *Api) SpanTrace(ctx *gin.Context) {
	// 从上下文中获取当前span
	span := trace.SpanFromContext(ctx.Request.Context())
	defer span.End()
	// 添加自定义属性到span
	span.SetAttributes(attribute.String("http.method", ctx.Request.Method))
	span.SetAttributes(attribute.String("handler.method", "span-trace"))
	// 记录自定义事件
	span.AddEvent("mock fetching user from database", trace.WithAttributes(
		attribute.String("user.id", ctx.Query("id")),
	))
	// 模拟数据库查询
	name, err := mockFetchUserNameFromDB(ctx.Request.Context(), ctx.Query("id"))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		response.FailWithError(ctx, err)
		return
	}
	// 添加更多属性
	span.SetAttributes(attribute.String("user.name", name))
	response.SuccessWithData[string](ctx, "trace")
}

// 模拟数据库操作函数,验证创建子span
func mockFetchUserNameFromDB(ctx context.Context, uid string) (string, error) {
	// 创建子span
	tracer := otel.Tracer("database")
	tracerCtx, span := tracer.Start(ctx, "mockFetchUserNameFromDB")
	defer span.End()
	if err := tracerCtx.Err(); err != nil {
		return "", err
	}
	// 添加属性
	span.SetAttributes(attribute.String("db.operation", "select"))
	span.SetAttributes(attribute.String("db.table", "users"))
	span.SetAttributes(attribute.String("user.id", uid))
	// 模拟数据库查询
	time.Sleep(10 * time.Millisecond)
	// 这里应该是实际的数据库查询代码
	if uid == "123" {
		return "weaver", nil
	}
	return "", response.UserNotExist
}
