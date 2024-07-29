package captcha

import (
	"context"
	"fmt"
	"gin-web/pkg/global"
	"github.com/spf13/viper"
	"time"
)

type RedisStore struct {
}

var ctx = context.Background()

func getKeyPrefix() string {
	return viper.GetString("captcha.keyPrefix")
}

func getExpiration() time.Duration {
	expiration := viper.GetDuration("captcha.expiration")
	return expiration * time.Second
}

func (r RedisStore) Set(id string, value string) error {
	return global.RedisClient.Client.Set(ctx, fmt.Sprintf("%s%s", getKeyPrefix(), id), value, getExpiration()).Err()
}

func (r RedisStore) Get(id string, clear bool) string {
	result, err := global.RedisClient.Client.Get(ctx, fmt.Sprintf("%s%s", getKeyPrefix(), id)).Result()
	if err != nil {
		return ""
	}
	if clear {
		delErr := global.RedisClient.Client.Del(ctx, fmt.Sprintf("%s%s", getKeyPrefix(), id)).Err()
		if delErr != nil {
			return ""
		}
	}
	return result
}

func (r RedisStore) Verify(id, answer string, clear bool) bool {
	result := r.Get(id, clear)
	return result == answer
}
