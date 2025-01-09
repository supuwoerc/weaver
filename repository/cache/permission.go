package cache

import (
	"sync"
)

type PermissionCache struct {
	*BasicCache
}

const PermissionCacheKey = "permission_cache"

var (
	permissionCache     *PermissionCache
	permissionCacheOnce sync.Once
)

func NewPermissionCache() *PermissionCache {
	permissionCacheOnce.Do(func() {
		permissionCache = &PermissionCache{BasicCache: NewBasicCache()}
	})
	return permissionCache
}
