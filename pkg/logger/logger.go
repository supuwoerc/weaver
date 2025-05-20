package logger

import (
	"context"

	"go.uber.org/zap"
)

type TraceContextKey string

const (
	TraceIdContextKey TraceContextKey = "trace_id"
)

type Logger struct {
	*zap.SugaredLogger
}

func NewLogger(z *zap.SugaredLogger) *Logger {
	return &Logger{
		SugaredLogger: z,
	}
}
func (l *Logger) WithContext(ctx context.Context) *zap.SugaredLogger {
	value := ctx.Value(string(TraceIdContextKey))
	result := l.SugaredLogger
	if value != nil {
		// generate new logger
		result = result.With(zap.String(string(TraceIdContextKey), value.(string)))
	}
	return result
}
