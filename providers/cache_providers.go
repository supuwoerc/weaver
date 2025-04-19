package providers

import (
	"gin-web/pkg/cache"
	"github.com/google/wire"
)

type Dept cache.SystemCache
type Permission cache.SystemCache

func SystemCaches(dept Dept, p Permission) []cache.SystemCache {
	return []cache.SystemCache{dept, p}
}

var SystemCacheProvider = wire.NewSet(
	SystemCaches,
	cache.NewSystemCacheManager,
)
