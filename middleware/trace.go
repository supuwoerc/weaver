package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/supuwoerc/weaver/conf"
	"github.com/supuwoerc/weaver/pkg/logger"
	"go.opentelemetry.io/otel/trace"
)

type TraceMiddleware struct {
	conf   *conf.Config
	logger *logger.Logger
}

func NewTraceMiddleware(conf *conf.Config, logger *logger.Logger) *TraceMiddleware {
	return &TraceMiddleware{
		conf:   conf,
		logger: logger,
	}
}

func (c *TraceMiddleware) Trace() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 设置 opentelemetry span context tract id 到 header
		requestTraceID := trace.SpanFromContext(ctx.Request.Context()).SpanContext().TraceID().String()
		ctx.Header(c.conf.System.TraceKey, requestTraceID)
		ctx.Set(string(logger.TraceIDContextKey), requestTraceID)
	}
}
