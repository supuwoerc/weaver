package providers

import (
	"github.com/supuwoerc/weaver/pkg/cache"

	"github.com/google/wire"
)

type DepartmentCache cache.SystemCache
type PermissionCache cache.SystemCache

func SystemCaches(dept DepartmentCache, p PermissionCache) []cache.SystemCache {
	return []cache.SystemCache{dept, p}
}

var SystemCacheProvider = wire.NewSet(
	SystemCaches,
	cache.NewSystemCacheManager,
)
