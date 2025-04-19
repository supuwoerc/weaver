package cache

import (
	"fmt"
	"gin-web/pkg/constant"
	"github.com/samber/lo"
)

type SystemCache interface {
	Key() constant.CacheKey
	Refresh() error
	Clean() error
}
type SystemCacheManager struct {
	Caches []SystemCache
}

func NewSystemCacheManager(caches ...SystemCache) *SystemCacheManager {
	return &SystemCacheManager{
		Caches: caches,
	}
}

func (s *SystemCacheManager) Refresh(keys ...constant.CacheKey) error {
	for _, key := range keys {
		cache, ok := lo.Find(s.Caches, func(item SystemCache) bool {
			return item.Key() == key
		})
		if ok {
			if e := cache.Refresh(); e != nil {
				return e
			}
		} else {
			return fmt.Errorf("refresh cache:cache %s not found", key)
		}
	}
	return nil
}

func (s *SystemCacheManager) Clean(keys ...constant.CacheKey) error {
	for _, key := range keys {
		cache, ok := lo.Find(s.Caches, func(item SystemCache) bool {
			return item.Key() == key
		})
		if ok {
			if e := cache.Clean(); e != nil {
				return e
			}
		} else {
			return fmt.Errorf("clean cache:cache %s not found", key)
		}
	}
	return nil
}
