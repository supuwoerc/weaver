package cache

import (
	"gin-web/pkg/redis"
	"sync"
)

type PermissionCache struct {
	redis *redis.CommonRedisClient
}

var (
	permissionCache     *PermissionCache
	permissionCacheOnce sync.Once
)

func NewPermissionCache(r *redis.CommonRedisClient) *PermissionCache {
	permissionCacheOnce.Do(func() {
		permissionCache = &PermissionCache{redis: r}
	})
	return permissionCache
}
