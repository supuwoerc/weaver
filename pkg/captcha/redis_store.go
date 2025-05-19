package captcha

import (
	"context"
	"fmt"
	"time"

	"github.com/supuwoerc/weaver/conf"
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/redis"
)

var ctx = context.Background()

type RedisStore struct {
	redisClient *redis.CommonRedisClient
	conf        *conf.Config
}

func NewRedisStore(r *redis.CommonRedisClient, conf *conf.Config) *RedisStore {
	return &RedisStore{
		redisClient: r,
		conf:        conf,
	}
}

func (r *RedisStore) getExpiration() time.Duration {
	expiration := r.conf.Captcha.Expiration
	return expiration * time.Second
}

func (r *RedisStore) Set(id string, value string) error {
	return r.redisClient.Client.Set(ctx, fmt.Sprintf("%s%s", constant.CaptchaCodePrefix, id), value, r.getExpiration()).Err()
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
