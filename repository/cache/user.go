package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/constant"
	"gin-web/pkg/response"
	"github.com/gin-gonic/gin"
)

type UserCache struct {
	*BasicCache
}

const USER_CACHE_KEY = "user_cache"

var (
	tokenPairKey = fmt.Sprintf("%s%s", USER_CACHE_KEY, "_token")
)

var userCache *UserCache

func NewUserCache(ctx *gin.Context) *UserCache {
	if userCache == nil {
		userCache = &UserCache{BasicCache: NewBasicCache(ctx)}
	}
	return userCache
}

func (u *UserCache) HSetTokenPair(ctx context.Context, email string, pair *models.TokenPair) error {
	if pair == nil {
		return constant.GetError(u.ctx, response.USER_LOGIN_TOKEN_PAIR_CACHE_ERR)
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
