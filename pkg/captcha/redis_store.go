package captcha

import (
	"context"
	"fmt"
	"gin-web/pkg/constant"
	"gin-web/pkg/redis"
	"github.com/spf13/viper"
	"sync"
	"time"
)

var ctx = context.Background()

type RedisStore struct {
	redisClient *redis.CommonRedisClient
}

var (
	redisStoreOnce sync.Once
	redisStore     *RedisStore
)

func NewRedisStore(r *redis.CommonRedisClient) *RedisStore {
	redisStoreOnce.Do(func() {
		redisStore = &RedisStore{redisClient: r}
	})
	return redisStore
}

func getExpiration() time.Duration {
	expiration := viper.GetDuration("captcha.expiration")
	return expiration * time.Second
}

func (r *RedisStore) Set(id string, value string) error {
	return r.redisClient.Client.Set(ctx, fmt.Sprintf("%s%s", constant.CaptchaCodePrefix, id), value, getExpiration()).Err()
}

func (r *RedisStore) Get(id string, clear bool) string {
	result, err := r.redisClient.Client.Get(ctx, fmt.Sprintf("%s%s", constant.CaptchaCodePrefix, id)).Result()
	if err != nil {
		return ""
	}
	if clear {
		delErr := r.redisClient.Client.Del(ctx, fmt.Sprintf("%s%s", constant.CaptchaCodePrefix, id)).Err()
		if delErr != nil {
			return ""
		}
	}
	return result
}

func (r *RedisStore) Verify(id, answer string, clear bool) bool {
	result := r.Get(id, clear)
	return result == answer
}
