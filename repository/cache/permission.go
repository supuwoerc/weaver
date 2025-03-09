package cache

import (
	"gin-web/pkg/redis"
)

type PermissionCache struct {
	redis *redis.CommonRedisClient
}

func NewPermissionCache(r *redis.CommonRedisClient) *PermissionCache {
	return &PermissionCache{redis: r}
}
