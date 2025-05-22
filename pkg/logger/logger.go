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

// WithContext Write the information in the context to a new logger and return
func (l *Logger) WithContext(ctx context.Context) *zap.SugaredLogger {
	keyString := string(TraceIDContextKey)
	key := TraceIDContextKey
	value := ctx.Value(key)
	result := l.SugaredLogger
	if value != nil {
		result = result.With(zap.String(keyString, value.(string)))
	} else {
		value = ctx.Value(keyString)
		if value != nil {
			result = result.With(zap.String(keyString, value.(string)))
		}
	}
	return result
}
