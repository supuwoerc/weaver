package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/supuwoerc/weaver/models"
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/redis"
	"github.com/supuwoerc/weaver/pkg/response"
)

type UserCache struct {
	redis *redis.CommonRedisClient
}

func NewUserCache(r *redis.CommonRedisClient) *UserCache {
	return &UserCache{redis: r}
}

func (u *UserCache) CacheTokenPair(ctx context.Context, email string, pair *models.TokenPair) error {
	if pair == nil {
		return response.UserLoginTokenPairCacheErr
	}
	return u.redis.Client.HSet(ctx, constant.UserTokenPairKey, email, pair).Err()
}

func (u *UserCache) GetTokenPairIsExist(ctx context.Context, email string) (bool, error) {
	return u.redis.Client.HExists(ctx, constant.UserTokenPairKey, email).Result()
}

func (u *UserCache) HDelTokenPair(ctx context.Context, email string) error {
	return u.redis.Client.HDel(ctx, constant.UserTokenPairKey, email).Err()
}

func (u *UserCache) GetTokenPair(ctx context.Context, email string) (*models.TokenPair, error) {
	var ret models.TokenPair
	err := u.redis.Client.HGet(ctx, constant.UserTokenPairKey, email).Scan(&ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (u *UserCache) activeAccountKey(id uint) string {
	return fmt.Sprintf("%s%d", constant.ActiveAccountPrefix, id)
}

func (u *UserCache) CacheActiveAccountCode(ctx context.Context, id uint, code string, duration time.Duration) error {
	return u.redis.Client.Set(ctx, u.activeAccountKey(id), code, duration).Err()
}

func (u *UserCache) GetActiveAccountCode(ctx context.Context, id uint) (string, error) {
	result, err := u.redis.Client.Get(ctx, u.activeAccountKey(id)).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (u *UserCache) RemoveActiveAccountCode(ctx context.Context, id uint) error {
	return u.redis.Client.Del(ctx, u.activeAccountKey(id)).Err()
}
