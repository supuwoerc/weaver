package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/supuwoerc/weaver/pkg/logger"
)

type EnginLoggerMiddleware struct {
	logger *logger.Logger
}

func NewEnginLoggerMiddleware(logger *logger.Logger) *EnginLoggerMiddleware {
	return &EnginLoggerMiddleware{
		logger: logger,
	}
}

func (e *EnginLoggerMiddleware) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		// Process request
		c.Next()
		param := gin.LogFormatterParams{
			Request: c.Request,
			Keys:    c.Keys,
		}
		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)
		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		param.BodySize = c.Writer.Size()
		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path

		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}
		e.logger.WithContext(c).Infow("http request",
			"path", param.Path,
			"latency", fmt.Sprintf("%v", param.Latency),
			"method", param.Method,
			"code", param.StatusCode,
			"error", param.ErrorMessage,
			"size", param.BodySize,
			"client", param.ClientIP,
		)
	}
}
