package cache

import (
	"context"
	"encoding/json"

	"github.com/supuwoerc/weaver/models"
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/redis"

	"github.com/samber/lo"
)

type DepartmentCache struct {
	redis *redis.CommonRedisClient
}

func NewDepartmentCache(r *redis.CommonRedisClient) *DepartmentCache {
	return &DepartmentCache{
		redis: r,
	}
}

func (d *DepartmentCache) CacheDepartment(ctx context.Context, key constant.CacheKey, depts []*models.Department) error {
	result, err := json.Marshal(depts)
	if err != nil {
		return err
	}
	return d.redis.Client.Set(ctx, string(key), string(result), 0).Err()
}

func (d *DepartmentCache) RemoveDepartmentCache(ctx context.Context, keys ...constant.CacheKey) error {
	return d.redis.Client.Del(ctx, lo.Map(keys, func(item constant.CacheKey, _ int) string {
		return string(item)
	})...).Err()
}

func (d *DepartmentCache) GetDepartmentCache(ctx context.Context, key constant.CacheKey) ([]*models.Department, error) {
	result, err := d.redis.Client.Get(ctx, string(key)).Result()
	if err != nil {
		return nil, err
	}
	var depts []*models.Department
	if err = json.Unmarshal([]byte(result), &depts); err != nil {
		return nil, err
	}
	return depts, nil
}
