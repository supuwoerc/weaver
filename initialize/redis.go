package initialize

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"github.com/supuwoerc/weaver/conf"
	weaverLogger "github.com/supuwoerc/weaver/pkg/logger"
	local "github.com/supuwoerc/weaver/pkg/redis"
)

type RedisLogLevel int

const (
	Silent RedisLogLevel = iota + 1
	Error
	Warn
	Info
)

type RedisLogger struct {
	goredislib.Hook
	*weaverLogger.Logger
	Level RedisLogLevel
}

func NewRedisLogger(l *weaverLogger.Logger, conf *conf.Config) *RedisLogger {
	return &RedisLogger{
		Logger: l,
		Level:  RedisLogLevel(conf.Logger.RedisLevel),
	}
}

func (r *RedisLogger) DialHook(next goredislib.DialHook) goredislib.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		if r.Level >= Info {
			r.Logger.WithContext(ctx).Infow("dialing to Redis", "network", network, "addr", addr)
		}
		conn, err := next(ctx, network, addr)
		if err != nil && r.Level >= Error {
			r.Logger.WithContext(ctx).Errorw("dialing Error", "error", err.Error())
		} else if r.Level >= Info {
			r.Logger.WithContext(ctx).Infow("successfully connected to Redis", "network", network, "addr", addr)
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
		if r.Level >= Info {
			r.Logger.WithContext(ctx).Infow("preparing to execute command", "command", cmd.Name())
		}
		err := next(ctx, cmd)
		if err != nil && r.Level >= Error {
			r.Logger.WithContext(ctx).Errorw("error executing command", "error", err.Error())
		} else if r.Level >= Info {
			r.Logger.WithContext(ctx).Infow("successfully executed command")
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
			if r.Level >= Info {
				r.Logger.WithContext(ctx).Infow("preparing to execute command in pipeline", "command", cmd.Name())
			}
		}
		err := next(ctx, cmds)
		if err != nil && r.Level >= Error {
			r.Logger.WithContext(ctx).Errorw("error executing commands in pipeline", "error", err.Error())
		} else if r.Level >= Info {
			r.Logger.WithContext(ctx).Infow("successfully executed commands in pipelined")
		}
		return err
	}
}

func NewRedisClient(hook goredislib.Hook, conf *conf.Config) *local.CommonRedisClient {
	client := goredislib.NewClient(&goredislib.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})
	client.AddHook(hook)
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	pool := goredis.NewPool(client)
	return local.NewCommonRedisClient(client, redsync.New(pool))
}
