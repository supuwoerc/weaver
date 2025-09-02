package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/redis"
)

type UserCache struct {
	redis *redis.CommonRedisClient
}

func NewUserCache(r *redis.CommonRedisClient) *UserCache {
	return &UserCache{redis: r}
}

func (u *UserCache) refreshTokenCacheKey(email string) string {
	return fmt.Sprintf("%s:%s", constant.UserRefreshTokenKey, email)
}

// CacheRefreshToken 存储用户的refreshToken
func (u *UserCache) CacheRefreshToken(ctx context.Context, email, refreshToken string, expiration time.Duration) error {
	return u.redis.Client.Set(ctx, u.refreshTokenCacheKey(email), refreshToken, expiration).Err()
}

// DeleteRefreshToken 删除用户的refreshToken
func (u *UserCache) DeleteRefreshToken(ctx context.Context, email string) error {
	return u.redis.Client.Del(ctx, u.refreshTokenCacheKey(email)).Err()
}

// GetRefreshToken 获取用户的refreshToken
func (u *UserCache) GetRefreshToken(ctx context.Context, email string) (string, error) {
	refreshToken, err := u.redis.Client.Get(ctx, u.refreshTokenCacheKey(email)).Result()
	if err != nil {
		return "", err
	}
	return refreshToken, nil
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
