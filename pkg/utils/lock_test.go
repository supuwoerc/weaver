package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supuwoerc/weaver/pkg/logger"
	"github.com/supuwoerc/weaver/pkg/redis"
	"go.uber.org/zap/zaptest"
)

func TestNewRedisLocksmith(t *testing.T) {
	t.Run("RedisLocksmith fields", func(t *testing.T) {
		l := logger.NewLogger(zaptest.NewLogger(t).Sugar())
		client := &redis.CommonRedisClient{}
		locksmith := NewRedisLocksmith(l, client)
		assert.NotNil(t, locksmith)
		assert.Equal(t, locksmith.logger, l)
		assert.Equal(t, locksmith.redisClient, client)
	})
}
