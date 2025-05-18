package redis

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
)

type CommonRedisClient struct {
	Client  redis.UniversalClient
	Redsync *redsync.Redsync
}

func NewCommonRedisClient(client redis.UniversalClient, redsync *redsync.Redsync) *CommonRedisClient {
	return &CommonRedisClient{
		Client:  client,
		Redsync: redsync,
	}
}

func (c *CommonRedisClient) NewMutex(name string, options ...redsync.Option) *redsync.Mutex {
	return c.Redsync.NewMutex(name, options...)
}
