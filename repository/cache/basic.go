package cache

import (
	"gin-web/pkg/redis"
	"sync"
)

var (
	basicCache     *BasicCache
	basicCacheOnce sync.Once
)

type BasicCache struct {
	redis *redis.CommonRedisClient
}

func NewBasicCache(r *redis.CommonRedisClient) *BasicCache {
	basicCacheOnce.Do(func() {
		basicCache = &BasicCache{
			redis: r,
		}
	})
	return basicCache
}
