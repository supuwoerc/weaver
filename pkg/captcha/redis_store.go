package captcha

import (
	"context"
	"fmt"
	"gin-web/pkg/constant"
	"gin-web/pkg/global"
	"github.com/spf13/viper"
	"time"
)

var ctx = context.Background()

type RedisStore struct {
}

func getExpiration() time.Duration {
	expiration := viper.GetDuration("captcha.expiration")
	return expiration * time.Second
}

func (r *RedisStore) Set(id string, value string) error {
	return global.RedisClient.Client.Set(ctx, fmt.Sprintf("%s%s", constant.CaptchaCodePrefix, id), value, getExpiration()).Err()
}

func (r *RedisStore) Get(id string, clear bool) string {
	result, err := global.RedisClient.Client.Get(ctx, fmt.Sprintf("%s%s", constant.CaptchaCodePrefix, id)).Result()
	if err != nil {
		return ""
	}
	if clear {
		delErr := global.RedisClient.Client.Del(ctx, fmt.Sprintf("%s%s", constant.CaptchaCodePrefix, id)).Err()
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
