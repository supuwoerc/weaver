package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

const DEFAULT_DURATION = 30 * 24 * 60 * 60 * time.Second

type RedisClient struct {
	*redis.Client
}

func (r *RedisClient) Get(key string) (any, error) {
	return r.Client.Get(context.Background(), key).Result()
}

func (r *RedisClient) Set(key string, value any, duration time.Duration) error {
	if duration <= 0 {
		duration = DEFAULT_DURATION
	}
	return r.Client.Set(context.Background(), key, value, duration).Err()
}

func (r *RedisClient) Del(key ...string) error {
	return r.Client.Del(context.Background(), key...).Err()
}

func (r *RedisClient) GetExpireDuration(key string) (time.Duration, error) {
	return r.Client.TTL(context.Background(), key).Result()
}
