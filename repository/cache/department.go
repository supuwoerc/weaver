package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/supuwoerc/weaver/models"
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/redis"

	"github.com/samber/lo"
)

type DepartmentCache struct {
	redis *redis.CommonRedisClient
	pid   int
}

func NewDepartmentCache(r *redis.CommonRedisClient) *DepartmentCache {
	return &DepartmentCache{
		redis: r,
		pid:   os.Getpid(),
	}
}

func (d *DepartmentCache) cacheKeyWithPid(key constant.CacheKey) string {
	return fmt.Sprintf("%s:%d", key, d.pid)
}

func (d *DepartmentCache) CacheDepartment(ctx context.Context, key constant.CacheKey, depts []*models.Department) error {
	result, err := json.Marshal(depts)
	if err != nil {
		return err
	}
	return d.redis.Client.Set(ctx, d.cacheKeyWithPid(key), string(result), 0).Err()
}

func (d *DepartmentCache) RemoveDepartmentCache(ctx context.Context, keys ...constant.CacheKey) error {
	return d.redis.Client.Del(ctx, lo.Map(keys, func(item constant.CacheKey, _ int) string {
		return d.cacheKeyWithPid(item)
	})...).Err()
}

func (d *DepartmentCache) GetDepartmentCache(ctx context.Context, key constant.CacheKey) ([]*models.Department, error) {
	result, err := d.redis.Client.Get(ctx, d.cacheKeyWithPid(key)).Result()
	if err != nil {
		return nil, err
	}
	var depts []*models.Department
	if err = json.Unmarshal([]byte(result), &depts); err != nil {
		return nil, err
	}
	return depts, nil
}
