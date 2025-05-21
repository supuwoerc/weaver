package utils

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/logger"
	"github.com/supuwoerc/weaver/pkg/redis"
	"github.com/supuwoerc/weaver/pkg/response"

	"github.com/go-redsync/redsync/v4"
	"github.com/pkg/errors"
)

type LocksmithLogger interface {
	logger.LogCtxInterface
	Errorw(msg string, keysAndValues ...interface{})
}

type LocksmithEmailClient interface {
	Alarm2Admin(ctx context.Context, subject constant.Subject, body string) error
}

type LocksmithMutex interface {
	NewMutex(name string, options ...redsync.Option) *redsync.Mutex
}

type RedisLocksmith struct {
	redisClient LocksmithMutex
	emailClient LocksmithEmailClient
	logger      LocksmithLogger
}

type RedisLock struct {
	*redsync.Mutex
	duration    time.Duration
	logger      LocksmithLogger
	emailClient LocksmithEmailClient
}

const (
	defaultMaxRetries   = 32
	defaultLockDuration = 10 * time.Second
)

func NewRedisLocksmith(logger LocksmithLogger, redisClient *redis.CommonRedisClient,
	emailClient LocksmithEmailClient) *RedisLocksmith {
	return &RedisLocksmith{
		logger:      logger,
		redisClient: redisClient,
		emailClient: emailClient,
	}
}

// NewLock 创建锁
func (r *RedisLocksmith) NewLock(t constant.Prefix, object ...string) *RedisLock {
	name := fmt.Sprintf("%s:%s", t, strings.Join(object, "_"))
	return &RedisLock{
		Mutex: r.redisClient.NewMutex(name,
			redsync.WithExpiry(defaultLockDuration),
			redsync.WithTries(defaultMaxRetries),
		),
		duration:    defaultLockDuration,
		logger:      r.logger,
		emailClient: r.emailClient,
	}
}

// Lock 会多次尝试(defaultMaxRetries 次), 如果尝试次数内还未获取到锁则返回错误
func (l *RedisLock) Lock(ctx context.Context, extend bool) error {
	err := l.Mutex.Lock()
	if err != nil {
		var e *redsync.ErrTaken
		if errors.As(err, &e) {
			return response.Busy
		}
		return fmt.Errorf("redis lock err: %v", err)
	}
	if extend {
		go l.autoExtend(ctx)
	}
	return nil
}

// TryLock 获取不到锁直接返回错误
func (l *RedisLock) TryLock(ctx context.Context, extend bool) error {
	err := l.Mutex.TryLock()
	if err != nil {
		var e *redsync.ErrTaken
		if errors.As(err, &e) {
			return response.Busy
		}
		return fmt.Errorf("redis lock err: %v", err)
	}
	if extend {
		go l.autoExtend(ctx)
	}
	return nil
}

// Unlock 解锁
func (l *RedisLock) Unlock() error {
	ok, err := l.Mutex.Unlock()
	if err != nil {
		if errors.Is(err, redsync.ErrLockAlreadyExpired) {
			return nil
		}
		return errors.Wrap(err, l.Name())
	}
	if !ok {
		return fmt.Errorf("%s unlock failed", l.Name())
	}
	return nil
}

func (l *RedisLock) alarm(ctx context.Context, subject constant.Subject, lockName string, err error) {
	l.logger.WithContext(ctx).Errorw("redis lock error alarm", "lock", lockName, "subject", subject, "err", err.Error())
	if err = l.emailClient.Alarm2Admin(
		ctx,
		subject,
		fmt.Sprintf("%s alarm: %v", lockName, err.Error()),
	); err != nil {
		l.logger.WithContext(ctx).Errorw("redis lock error alarm to admin err", "err", err.Error())
	}
}

func (l *RedisLock) unlockWithAlarm(ctx context.Context) {
	err := l.Unlock()
	if err != nil {
		go l.alarm(ctx, constant.UnlockFail, l.Name(), err)
	}
}

func (l *RedisLock) extendWithAlarm(ctx context.Context) bool {
	ok, err := l.Extend()
	if err != nil {
		go l.alarm(ctx, constant.ExtendErr, l.Name(), err)
		return false
	} else if !ok {
		go l.alarm(ctx, constant.ExtendFail, l.Name(), errors.New("lock couldn't be extended"))
		return false
	}
	return true
}

// autoExtend 自动延长锁的过期时间
func (l *RedisLock) autoExtend(ctx context.Context) {
	ticker := time.NewTicker(l.duration / 2)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			l.unlockWithAlarm(ctx)
			return
		case <-ticker.C:
			select {
			case <-ctx.Done():
				ticker.Stop()
				l.unlockWithAlarm(ctx)
				return
			default:
				if success := l.extendWithAlarm(ctx); !success {
					ticker.Stop()
					return
				}
			}
		}
	}
}
