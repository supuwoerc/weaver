package initialize

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/supuwoerc/weaver/conf"
	local "github.com/supuwoerc/weaver/pkg/redis"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
)

type RedisLogSyncer io.Writer
type RedisLogger struct {
	goredislib.Hook
	syncer RedisLogSyncer
}

func NewRedisLogger(syncer RedisLogSyncer) *RedisLogger {
	return &RedisLogger{
		syncer: syncer,
	}
}

func (r *RedisLogger) DialHook(next goredislib.DialHook) goredislib.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		_, _ = fmt.Fprintf(r.syncer, "{\"caller\":Redis,\"event\":Dialing to Redis,\"time\":\"%s\",\"network\":\"%s\",\"address\":\"%s\"}\n", time.Now().Format(time.DateTime), network, addr)
		conn, err := next(ctx, network, addr)
		if err != nil {
			_, _ = fmt.Fprintf(r.syncer, "{\"caller\":Redis,\"event\":Dialing Error,\"time\":\"%s\",\"error\":\"%s\"}\n", time.Now().Format(time.DateTime), err.Error())
		} else {
			_, _ = fmt.Fprintf(r.syncer, "{\"caller\":Redis,\"event\":Successfully connected to Redis,\"time\":\"%s\",\"network\":\"%s\",\"address\":\"%s\"}\n", time.Now().Format(time.DateTime), network, addr)
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
		_, _ = fmt.Fprintf(r.syncer, "{\"caller\":Redis,\"event\":Preparing to execute command,\"time\":\"%s\",\"command\":\"%s\",\"args\":\"%s\"}\n", time.Now().Format(time.DateTime), cmd.Name(), builder.String())
		err := next(ctx, cmd)
		if err != nil {
			_, _ = fmt.Fprintf(r.syncer, "{\"caller\":Redis,\"event\":Error executing command,\"time\":\"%s\",\"command\":\"%s\",\"args\":\"%s\"}\n", time.Now().Format(time.DateTime), cmd.Name(), err.Error())
		} else {
			_, _ = fmt.Fprintf(r.syncer, "{\"caller\":Redis,\"event\":Successfully executed command,\"time\":\"%s\",\"command\":\"%s\"}\n", time.Now().Format(time.DateTime), cmd.Name())
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
			_, _ = fmt.Fprintf(r.syncer, "{\"caller\":Redis,\"event\":Preparing to execute command in pipeline,\"time\":\"%s\",\"command\":\"%s\",\"args\":\"%s\"}\n", time.Now().Format(time.DateTime), cmd.Name(), builder.String())
		}
		err := next(ctx, cmds)
		if err != nil {
			_, _ = fmt.Fprintf(r.syncer, "{\"caller\":Redis,\"event\":Error executing commands in pipeline,\"time\":\"%s\",\"error\":\"%s\"}\n", time.Now().Format(time.DateTime), err.Error())
		} else {
			_, _ = fmt.Fprintf(r.syncer, "{\"caller\":Redis,\"event\":Successfully executed commands in pipeline,\"time\":\"%s\"}\n", time.Now().Format(time.DateTime))
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
