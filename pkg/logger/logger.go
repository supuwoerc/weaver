package logger

import (
	"context"

	"go.uber.org/zap"
)

type TraceContextKey string

const (
	TraceIDContextKey TraceContextKey = "trace_id"
)

type LogCtxInterface interface {
	WithContext(ctx context.Context) *zap.SugaredLogger
}

type Logger struct {
	LogCtxInterface
	*zap.SugaredLogger
}

func NewLogger(z *zap.SugaredLogger) *Logger {
	return &Logger{
		SugaredLogger: z,
	}
}
func (l *Logger) WithContext(ctx context.Context) *zap.SugaredLogger {
	value := ctx.Value(string(TraceIDContextKey))
	result := l.SugaredLogger
	if value != nil {
		// generate new logger
		result = result.With(zap.String(string(TraceIDContextKey), value.(string)))
	}
	return result
}
