package logger

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestNewLogger(t *testing.T) {
	zapLogger := zaptest.NewLogger(t).Sugar()
	logger := NewLogger(zapLogger)
	assert.NotNil(t, logger)
	assert.Equal(t, zapLogger, logger.SugaredLogger)
}

func TestLogger_WithContext(t *testing.T) {
	mockLogger := NewLogger(zaptest.NewLogger(t).Sugar())
	testTraceID := uuid.New().String()
	t.Run("context with traceID", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), TraceIDContextKey, testTraceID)
		logger := mockLogger.WithContext(ctx)
		assert.NotNil(t, logger)
		assert.NotEqual(t, logger, mockLogger)
	})

	t.Run("gin context with traceID", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(nil)
		ctx.Set(string(TraceIDContextKey), testTraceID)
		logger := mockLogger.WithContext(ctx)
		assert.NotNil(t, logger)
		assert.NotEqual(t, logger, mockLogger)
	})

	t.Run("context without traceID", func(t *testing.T) {
		ctx := context.Background()
		logger := mockLogger.WithContext(ctx)
		assert.Equal(t, mockLogger.SugaredLogger, logger)
	})
}
