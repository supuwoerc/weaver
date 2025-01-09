package cache

import (
	"gin-web/pkg/global"
	"gin-web/pkg/redis"
	"sync"
)

var (
	basicCache     *BasicCache
	basicCacheOnce sync.Once
)

type BasicCache struct {
	redis *redis.RedisClient
}

func NewBasicCache() *BasicCache {
	basicCacheOnce.Do(func() {
		basicCache = &BasicCache{
			redis: global.RedisClient,
		}
	})
	return basicCache
}
