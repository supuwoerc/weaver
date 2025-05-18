package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-redsync/redsync/v4"

	"github.com/go-redis/redismock/v9"
)

func TestNewCommonRedisClient(t *testing.T) {
	redisClientMock, _ := redismock.NewClientMock()
	redsyncMock := &redsync.Redsync{}
	commonRedisClient := NewCommonRedisClient(redisClientMock, redsyncMock)
	assert.Equal(t, commonRedisClient.Client, redisClientMock)
	assert.Equal(t, commonRedisClient.Redsync, redsyncMock)
}

func TestCommonRedisClient_NewMutex(t *testing.T) {
	redisClientMock, _ := redismock.NewClientMock()
	redsyncMock := &redsync.Redsync{}
	commonRedisClient := NewCommonRedisClient(redisClientMock, redsyncMock)
	assert.Equal(t, commonRedisClient.Client, redisClientMock)
	assert.Equal(t, commonRedisClient.Redsync, redsyncMock)
	name := "mutexName"
	mutex := commonRedisClient.NewMutex(name)
	assert.Equal(t, name, mutex.Name())
}
