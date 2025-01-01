package cache

import (
	"gin-web/pkg/global"
	"gin-web/pkg/redis"
)

type BasicCache struct {
	redis *redis.RedisClient
}

func NewBasicCache() *BasicCache {
	return &BasicCache{
		redis: global.RedisClient,
	}
}
