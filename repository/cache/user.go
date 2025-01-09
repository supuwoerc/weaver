package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/response"
	"sync"
)

type UserCache struct {
	*BasicCache
}

const UserCacheKey = "user_cache"

var (
	tokenPairKey  = fmt.Sprintf("%s%s", UserCacheKey, "_token")
	userCache     *UserCache
	userCacheOnce sync.Once
)

func NewUserCache() *UserCache {
	userCacheOnce.Do(func() {
		userCache = &UserCache{BasicCache: NewBasicCache()}
	})
	return userCache
}

func (u *UserCache) HSetTokenPair(ctx context.Context, email string, pair *models.TokenPair) error {
	if pair == nil {
		return response.UserLoginTokenPairCacheErr
	}
	result, err := json.Marshal(pair)
	if err != nil {
		return err
	}
	return u.redis.Client.HSet(ctx, tokenPairKey, email, result).Err()
}

func (u *UserCache) HExistsTokenPair(ctx context.Context, email string) (bool, error) {
	return u.redis.Client.HExists(ctx, tokenPairKey, email).Result()
}

func (u *UserCache) HDelTokenPair(ctx context.Context, email string) error {
	return u.redis.Client.HDel(ctx, tokenPairKey, email).Err()
}

func (u *UserCache) HGetTokenPair(ctx context.Context, email string) (*models.TokenPair, error) {
	result, err := u.redis.Client.HGet(ctx, tokenPairKey, email).Result()
	if err != nil {
		return nil, err
	}
	var ret models.TokenPair
	err = json.Unmarshal([]byte(result), &ret)
	return &ret, err
}
