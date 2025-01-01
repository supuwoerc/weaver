package initialize

import (
	"context"
	local "gin-web/pkg/redis"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() *local.RedisClient {
	client := goredislib.NewClient(&goredislib.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	pool := goredis.NewPool(client)
	return &local.RedisClient{Client: client, Redsync: redsync.New(pool)}
}
