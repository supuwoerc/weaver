package initialize

import (
	"context"
	wrapRedis "gin-web/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() *wrapRedis.RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"), // no password set
		DB:       viper.GetInt("redis.db"),          // use default DB
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return &wrapRedis.RedisClient{Client: client}
}
