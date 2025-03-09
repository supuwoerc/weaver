package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/constant"
	"gin-web/pkg/redis"
	"gin-web/pkg/response"
	"time"
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
	result, err := json.Marshal(pair)
	if err != nil {
		return err
	}
	return u.redis.Client.HSet(ctx, constant.UserTokenPairKey, email, result).Err()
}

func (u *UserCache) GetTokenPairIsExist(ctx context.Context, email string) (bool, error) {
	return u.redis.Client.HExists(ctx, constant.UserTokenPairKey, email).Result()
}

func (u *UserCache) HDelTokenPair(ctx context.Context, email string) error {
	return u.redis.Client.HDel(ctx, constant.UserTokenPairKey, email).Err()
}

func (u *UserCache) GetTokenPair(ctx context.Context, email string) (*models.TokenPair, error) {
	result, err := u.redis.Client.HGet(ctx, constant.UserTokenPairKey, email).Result()
	if err != nil {
		return nil, err
	}
	var ret models.TokenPair
	err = json.Unmarshal([]byte(result), &ret)
	return &ret, err
}

func (u *UserCache) CacheActiveAccountCode(ctx context.Context, id uint, code string, duration time.Duration) error {
	return u.redis.Client.Set(ctx, fmt.Sprintf("%s%d", constant.ActiveAccountPrefix, id), code, duration).Err()
}

func (u *UserCache) GetActiveAccountCode(ctx context.Context, id uint) (string, error) {
	result, err := u.redis.Client.Get(ctx, fmt.Sprintf("%s%d", constant.ActiveAccountPrefix, id)).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (u *UserCache) RemoveActiveAccountCode(ctx context.Context, id uint) error {
	return u.redis.Client.Del(ctx, fmt.Sprintf("%s%d", constant.ActiveAccountPrefix, id)).Err()
}
