package initialize

import (
	"context"
	"testing"

	"github.com/supuwoerc/weaver/conf"
	"github.com/supuwoerc/weaver/pkg/logger"
)

func BenchmarkLogger(b *testing.B) {
	setting := &conf.Config{
		Logger: conf.LoggerConfig{
			MaxSize:    100,
			MaxBackups: 10,
			MaxAge:     10,
			Level:      0,
			Dir:        "./log",
			Stdout:     true,
		},
	}
	zapLogger := NewZapLogger(setting, NewWriterSyncer(setting))
	newLogger := logger.NewLogger(zapLogger)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		newLogger.WithContext(context.Background()).Info("test")
	}
}
