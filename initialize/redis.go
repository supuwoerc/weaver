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
	"strings"
	"time"
)

type RedisLogger struct {
	logger io.Writer
}

func (r *RedisLogger) DialHook(next goredislib.DialHook) goredislib.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		_, _ = fmt.Fprint(r.logger, fmt.Sprintf("[Redis] [%s] Dialing to Redis at %s://%s\n", time.Now().Format(time.DateTime), network, addr))
		conn, err := next(ctx, network, addr)
		if err != nil {
			_, _ = fmt.Fprint(r.logger, fmt.Sprintf("[Redis] [%s] Error dialing Redis: %s\n", time.Now().Format(time.DateTime), err.Error()))
		} else {
			_, _ = fmt.Fprint(r.logger, fmt.Sprintf("[Redis] [%s] Successfully connected to Redis at %s://%s\n", time.Now().Format(time.DateTime), network, addr))
		}
		return conn, err
	}
}

func (r *RedisLogger) ProcessHook(next goredislib.ProcessHook) goredislib.ProcessHook {
	return func(ctx context.Context, cmd goredislib.Cmder) error {
		builder := strings.Builder{}
		for i, arg := range cmd.Args() {
			if i > 0 {
				builder.WriteString(" ")
			}
			builder.WriteString(fmt.Sprintf("%v", arg))
		}
		_, _ = fmt.Fprint(r.logger, fmt.Sprintf("[Redis] [%s] Preparing to execute command: %s, [Args]: %s\n", time.Now().Format(time.DateTime), cmd.Name(), builder.String()))
		err := next(ctx, cmd)
		if err != nil {
			_, _ = fmt.Fprint(r.logger, fmt.Sprintf("[Redis] [%s] Error executing command %s: %s\n", time.Now().Format(time.DateTime), cmd.Name(), err.Error()))
		} else {
			_, _ = fmt.Fprint(r.logger, fmt.Sprintf("[Redis] [%s] Successfully executed command: %s\n", time.Now().Format(time.DateTime), cmd.Name()))
		}
		return err
	}
}

func (r *RedisLogger) ProcessPipelineHook(next goredislib.ProcessPipelineHook) goredislib.ProcessPipelineHook {
	return func(ctx context.Context, cmds []goredislib.Cmder) error {
		for _, cmd := range cmds {
			builder := strings.Builder{}
			for i, arg := range cmd.Args() {
				if i > 0 {
					builder.WriteString(" ")
				}
				builder.WriteString(fmt.Sprintf("%v", arg))
			}
			_, _ = fmt.Fprint(r.logger, fmt.Sprintf("[Redis] [%s] Preparing to execute command in pipeline: %s, [Args]: %s\n", time.Now().Format(time.DateTime), cmd.Name(), builder.String()))
		}
		err := next(ctx, cmds)
		if err != nil {
			_, _ = fmt.Fprint(r.logger, fmt.Sprintf("[Redis] [%s] Error executing commands in pipeline: %s\n", time.Now().Format(time.DateTime), err.Error()))
		} else {
			_, _ = fmt.Fprint(r.logger, fmt.Sprintf("[Redis] [%s] Successfully executed commands in pipeline", time.Now().Format(time.DateTime)))
		}
		return err
	}
}

func InitRedis(logger io.Writer) *local.CommonRedisClient {
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
	return &local.CommonRedisClient{Client: client, Redsync: redsync.New(pool)}
}
