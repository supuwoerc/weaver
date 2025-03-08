package cache

import (
	"gin-web/pkg/redis"
	"sync"
)

type PermissionCache struct {
	*BasicCache
}

var (
	permissionCache     *PermissionCache
	permissionCacheOnce sync.Once
)

func NewPermissionCache(r *redis.CommonRedisClient) *PermissionCache {
	permissionCacheOnce.Do(func() {
		permissionCache = &PermissionCache{BasicCache: NewBasicCache(r)}
	})
	return permissionCache
}
