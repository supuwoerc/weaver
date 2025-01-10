package initialize

import (
	"context"
	local "gin-web/pkg/redis"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net"
)

type RedisLogger struct {
	logger *zap.SugaredLogger
}

func (r *RedisLogger) DialHook(next goredislib.DialHook) goredislib.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		r.logger.Infof("[Redis] Dialing to Redis at %s://%s", network, addr)
		conn, err := next(ctx, network, addr)
		if err != nil {
			r.logger.Errorf("[Redis] Error dialing Redis: %s", err.Error())
		} else {
			r.logger.Infof("[Redis] Successfully connected to Redis at %s://%s", network, addr)
		}
		return conn, err
	}
}

func (r *RedisLogger) ProcessHook(next goredislib.ProcessHook) goredislib.ProcessHook {
	return func(ctx context.Context, cmd goredislib.Cmder) error {
		r.logger.Infof("[Redis] Preparing to execute command: %s, Args: %s", cmd.Name(), cmd.Args())
		err := next(ctx, cmd)
		if err != nil {
			r.logger.Errorf("[Redis] Error executing command %s: %s", cmd.Name(), err.Error())
		} else {
			r.logger.Infof("[Redis] Successfully executed command: %s", cmd.Name())
		}
		return err
	}
}

func (r *RedisLogger) ProcessPipelineHook(next goredislib.ProcessPipelineHook) goredislib.ProcessPipelineHook {
	return func(ctx context.Context, cmds []goredislib.Cmder) error {
		for _, cmd := range cmds {
			r.logger.Infof("[Redis] Preparing to execute command in pipeline: %s, Args: %s", cmd.Name(), cmd.Args())
		}
		err := next(ctx, cmds)
		if err != nil {
			r.logger.Errorf("[Redis] Error executing commands in pipeline: %s", err.Error())
		} else {
			r.logger.Infof("[Redis] Successfully executed commands in pipeline")
		}
		return err
	}
}

func InitRedis(logger *zap.SugaredLogger) *local.RedisClient {
	client := goredislib.NewClient(&goredislib.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
	client.AddHook(&RedisLogger{
		logger: logger,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	pool := goredis.NewPool(client)
	return &local.RedisClient{Client: client, Redsync: redsync.New(pool)}
}
