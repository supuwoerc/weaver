package initialize

import (
	"context"
	"net"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/supuwoerc/weaver/conf"
	weaverLogger "github.com/supuwoerc/weaver/pkg/logger"
	"go.uber.org/zap/zaptest"
)

func TestNewRedisLogger(t *testing.T) {
	// 创建测试用的 logger
	logger := weaverLogger.NewLogger(zaptest.NewLogger(t).Sugar())
	t.Run("test logger with different levels", func(t *testing.T) {
		testCases := []struct {
			name     string
			logLevel RedisLogLevel
			logger   *weaverLogger.Logger
		}{
			{"silent level", Silent, logger},
			{"error level", Error, logger},
			{"warn level", Warn, logger},
			{"info level", Info, logger},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				config := &conf.Config{
					Logger: conf.LoggerConfig{
						RedisLevel: int(tc.logLevel),
					},
				}
				redisLogger := NewRedisLogger(logger, config)
				require.NotNil(t, redisLogger)
				require.Equal(t, tc.logger, redisLogger.Logger)
				require.Equal(t, tc.logLevel, redisLogger.Level)
			})
		}
	})
}

func TestRedisLogger_DialHook(t *testing.T) {
	logger := weaverLogger.NewLogger(zaptest.NewLogger(t).Sugar())
	config := &conf.Config{
		Logger: conf.LoggerConfig{
			RedisLevel: int(Info),
		},
	}
	redisLogger := NewRedisLogger(logger, config)
	t.Run("successful dial", func(t *testing.T) {
		ctx := context.Background()
		nextHook := func(ctx context.Context, network, addr string) (net.Conn, error) {
			return &net.TCPConn{}, nil
		}
		hook := redisLogger.DialHook(nextHook)
		conn, err := hook(ctx, "tcp", "localhost:6379")
		require.NoError(t, err)
		require.NotNil(t, conn)
	})
	t.Run("dial error", func(t *testing.T) {
		ctx := context.Background()
		nextHook := func(ctx context.Context, network, addr string) (net.Conn, error) {
			return nil, redis.ErrClosed
		}
		hook := redisLogger.DialHook(nextHook)
		conn, err := hook(ctx, "tcp", "localhost:6379")
		require.Error(t, err)
		require.Nil(t, conn)
	})
}

func TestRedisLogger_ProcessHook(t *testing.T) {
	logger := weaverLogger.NewLogger(zaptest.NewLogger(t).Sugar())
	config := &conf.Config{
		Logger: conf.LoggerConfig{
			RedisLevel: int(Info),
		},
	}
	redisLogger := NewRedisLogger(logger, config)
	t.Run("successful process", func(t *testing.T) {
		ctx := context.Background()
		cmd := redis.NewStringCmd(ctx, "GET", "key")
		nextHook := func(ctx context.Context, cmd redis.Cmder) error {
			return nil
		}
		hook := redisLogger.ProcessHook(nextHook)
		err := hook(ctx, cmd)
		require.NoError(t, err)
	})
	t.Run("process error", func(t *testing.T) {
		ctx := context.Background()
		cmd := redis.NewStringCmd(ctx, "GET", "key")
		nextHook := func(ctx context.Context, cmd redis.Cmder) error {
			return redis.ErrClosed
		}
		hook := redisLogger.ProcessHook(nextHook)
		err := hook(ctx, cmd)
		require.Error(t, err)
	})
}

func TestRedisLogger_ProcessPipelineHook(t *testing.T) {
	logger := weaverLogger.NewLogger(zaptest.NewLogger(t).Sugar())
	config := &conf.Config{
		Logger: conf.LoggerConfig{
			RedisLevel: int(Info),
		},
	}
	redisLogger := NewRedisLogger(logger, config)

	t.Run("successful pipeline", func(t *testing.T) {
		ctx := context.Background()
		cmds := []redis.Cmder{
			redis.NewStringCmd(ctx, "GET", "key1"),
			redis.NewStringCmd(ctx, "GET", "key2"),
		}
		nextHook := func(ctx context.Context, cmds []redis.Cmder) error {
			return nil
		}
		hook := redisLogger.ProcessPipelineHook(nextHook)
		err := hook(ctx, cmds)
		require.NoError(t, err)
	})

	t.Run("pipeline error", func(t *testing.T) {
		ctx := context.Background()
		cmds := []redis.Cmder{
			redis.NewStringCmd(ctx, "GET", "key1"),
			redis.NewStringCmd(ctx, "GET", "key2"),
		}
		nextHook := func(ctx context.Context, cmds []redis.Cmder) error {
			return redis.ErrClosed
		}
		hook := redisLogger.ProcessPipelineHook(nextHook)
		err := hook(ctx, cmds)
		require.Error(t, err)
	})
}
