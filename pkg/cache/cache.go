package cache

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/samber/lo"
)

type SystemCache interface {
	CacheKey() string
	RefreshCache(ctx context.Context) error
	CleanCache(ctx context.Context) error
}

//go:generate stringer -type=cacheOperate -linecomment -output cache_operate_string.go
type cacheOperate int

const (
	refresh cacheOperate = iota + 1 // refresh
	clean                           // clean
)

type SystemCacheManager struct {
	caches []SystemCache
}

func NewSystemCacheManager(caches ...SystemCache) *SystemCacheManager {
	return &SystemCacheManager{
		caches: lo.Filter(caches, func(item SystemCache, _ int) bool {
			return item != nil
		}),
	}
}

func (s *SystemCacheManager) Refresh(ctx context.Context, keys ...string) error {
	return operateCache(ctx, refresh, s, keys...)
}

func (s *SystemCacheManager) Clean(ctx context.Context, keys ...string) error {
	return operateCache(ctx, clean, s, keys...)
}

func operateCache(ctx context.Context, op cacheOperate, s *SystemCacheManager, keys ...string) error {
	if len(s.caches) == 0 {
		return fmt.Errorf("%s cache fail: cache slice is empty", op)
	}
	for _, key := range keys {
		cache, exists := lo.Find(s.caches, func(c SystemCache) bool {
			return c.CacheKey() == key
		})
		if !exists {
			return fmt.Errorf("%s cache fail: cache %s not found", op, key)
		}
		var err error
		switch op {
		case refresh:
			err = cache.RefreshCache(ctx)
		case clean:
			err = cache.CleanCache(ctx)
		default:
			return fmt.Errorf("%s is invalid operate", op)
		}
		if err != nil {
			return errors.Wrapf(err, "%s cache fail,key %s", op, key)
		}
	}
	return nil
}
