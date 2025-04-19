package cache

import (
	"context"
	"fmt"
	"github.com/samber/lo"
)

type SystemCache interface {
	Key() string
	Refresh(ctx context.Context) error
	Clean(ctx context.Context) error
}
type SystemCacheManager struct {
	Caches []SystemCache
}

func NewSystemCacheManager(caches ...SystemCache) *SystemCacheManager {
	return &SystemCacheManager{
		Caches: caches,
	}
}

func (s *SystemCacheManager) Refresh(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		cache, ok := lo.Find(s.Caches, func(item SystemCache) bool {
			return item.Key() == key
		})
		if ok {
			if e := cache.Refresh(ctx); e != nil {
				return e
			}
		} else {
			return fmt.Errorf("refresh cache fail:cache %s not found", key)
		}
	}
	return nil
}

func (s *SystemCacheManager) Clean(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		cache, ok := lo.Find(s.Caches, func(item SystemCache) bool {
			return item.Key() == key
		})
		if ok {
			if e := cache.Clean(ctx); e != nil {
				return e
			}
		} else {
			return fmt.Errorf("clean cache fail:cache %s not found", key)
		}
	}
	return nil
}
