package redis

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client  *redis.Client
	Redsync *redsync.Redsync
}
