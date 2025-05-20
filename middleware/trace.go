package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/supuwoerc/weaver/conf"
	"go.uber.org/zap"
)

type TraceMiddleware struct {
	conf   *conf.Config
	logger *zap.SugaredLogger
}

type TraceContextKey string

const (
	TraceIdContextKey TraceContextKey = "trace_id"
)

func NewTraceMiddleware(conf *conf.Config, logger *zap.SugaredLogger) *TraceMiddleware {
	return &TraceMiddleware{
		conf:   conf,
		logger: logger,
	}
}

func (c *TraceMiddleware) Trace() gin.HandlerFunc {
	return func(context *gin.Context) {
		requestTraceId := context.GetHeader(c.conf.System.TraceKey)
		if strings.TrimSpace(requestTraceId) == "" {
			requestTraceId = c.generateTraceId()
		}
		context.Set(string(TraceIdContextKey), requestTraceId)
	}
}

func (c *TraceMiddleware) generateTraceId() string {
	return uuid.New().String()
}
