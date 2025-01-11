package initialize

import (
	"context"
	"fmt"
	local "gin-web/pkg/redis"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"io"
	"net"
)

type RedisLogger struct {
	logger io.Writer
}

func (r *RedisLogger) DialHook(next goredislib.DialHook) goredislib.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		_, _ = fmt.Fprint(r.logger, fmt.Sprintf("[Redis] Dialing to Redis at %s://%s\n", network, addr))
		conn, err := next(ctx, network, addr)
		if err != nil {
			_, _ = fmt.Fprint(r.logger, fmt.Sprintf("[Redis] Error dialing Redis: %s\n", err.Error()))
		} else {
			_, _ = fmt.Fprint(r.logger, fmt.Sprintf("[Redis] Successfully connected to Redis at %s://%s\n", network, addr))
		}
		return conn, err
	}
}

func (r *RedisLogger) ProcessHook(next goredislib.ProcessHook) goredislib.ProcessHook {
	return func(ctx context.Context, cmd goredislib.Cmder) error {
		_, _ = fmt.Fprint(r.logger, fmt.Sprintf("[Redis] Preparing to execute command: %s, Args: %s\n", cmd.Name(), cmd.Args()))
		err := next(ctx, cmd)
		if err != nil {
			_, _ = fmt.Fprint(r.logger, fmt.Sprintf("[Redis] Error executing command %s: %s\n", cmd.Name(), err.Error()))
		} else {
			_, _ = fmt.Fprint(r.logger, fmt.Sprintf("[Redis] Successfully executed command: %s\n", cmd.Name()))
		}
		return err
	}
}

func (r *RedisLogger) ProcessPipelineHook(next goredislib.ProcessPipelineHook) goredislib.ProcessPipelineHook {
	return func(ctx context.Context, cmds []goredislib.Cmder) error {
		for _, cmd := range cmds {
			_, _ = fmt.Fprint(r.logger, fmt.Sprintf("[Redis] Preparing to execute command in pipeline: %s, Args: %s\n", cmd.Name(), cmd.Args()))
		}
		err := next(ctx, cmds)
		if err != nil {
			_, _ = fmt.Fprint(r.logger, fmt.Sprintf("[Redis] Error executing commands in pipeline: %s\n", err.Error()))
		} else {
			_, _ = fmt.Fprint(r.logger, fmt.Sprintf("[Redis] Successfully executed commands in pipeline"))
		}
		return err
	}
}

func InitRedis(logger io.Writer) *local.RedisClient {
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
