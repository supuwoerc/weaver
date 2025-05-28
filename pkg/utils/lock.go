package utils

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
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
}

type LocksmithMutex interface {
	NewMutex(name string, options ...redsync.Option) *redsync.Mutex
}

type RedisLocksmith struct {
	redisClient LocksmithMutex
	logger      LocksmithLogger
}

const (
	lockStateUnlocked = 0 // 未锁定
	lockStateLocked   = 1 // 锁定
	lockStateReleased = 2 // 释放
)

type RedisLock struct {
	*redsync.Mutex
	duration time.Duration
	logger   LocksmithLogger
	state    atomic.Int64
}

const (
	defaultMaxRetries   = 32
	defaultLockDuration = 10 * time.Second
)

func NewRedisLocksmith(logger LocksmithLogger, redisClient *redis.CommonRedisClient) *RedisLocksmith {
	return &RedisLocksmith{
		logger:      logger,
		redisClient: redisClient,
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
		duration: defaultLockDuration,
		logger:   r.logger,
	}
}

// Lock 会多次尝试(defaultMaxRetries 次), 如果尝试次数内还未获取到锁则返回错误
func (l *RedisLock) Lock(ctx context.Context, extend bool) error {
	return l.acquire(ctx, l.Mutex.Lock, extend)
}

// TryLock 获取不到锁直接返回错误
func (l *RedisLock) TryLock(ctx context.Context, extend bool) error {
	return l.acquire(ctx, l.Mutex.TryLock, extend)
}

func (l *RedisLock) acquire(ctx context.Context, lockMethod func() error, extend bool) error {
	if !l.state.CompareAndSwap(lockStateUnlocked, lockStateLocked) {
		currentState := l.state.Load()
		switch currentState {
		case lockStateLocked:
			return response.Busy
		case lockStateReleased:
			return fmt.Errorf("lock has been released")
		default:
			return fmt.Errorf("lock is in invalid state: %d", currentState)
		}
	}
	err := lockMethod()
	if err != nil {
		l.state.Store(lockStateUnlocked)
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
	if !l.state.CompareAndSwap(lockStateLocked, lockStateReleased) {
		currentState := l.state.Load()
		switch currentState {
		case lockStateUnlocked:
			return fmt.Errorf("cannot unlock: lock not held")
		case lockStateReleased:
			return nil
		default:
			return fmt.Errorf("cannot unlock: lock in invalid state: %d", currentState)
		}
	}
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

func (l *RedisLock) unlockWithLog(ctx context.Context) {
	err := l.Unlock()
	if err != nil {
		l.logger.WithContext(ctx).Errorw("redis unlock fail", "err", err.Error())
	}
}

func (l *RedisLock) extend(ctx context.Context) bool {
	if l.state.Load() != lockStateLocked {
		return false
	}
	ok, err := l.Mutex.Extend()
	if err != nil {
		l.logger.WithContext(ctx).Errorw("lock couldn't be extended", "err", err.Error(), "name", l.Name())
		return false
	} else if !ok {
		l.logger.WithContext(ctx).Errorw("redis extend lock fail ", "name", l.Name())
		return false
	}
	return true
}

// autoExtend 自动延长锁的过期时间
func (l *RedisLock) autoExtend(ctx context.Context) {
	ticker := time.NewTicker(l.duration / 3)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			l.unlockWithLog(ctx)
			return
		case <-ticker.C:
			select {
			case <-ctx.Done():
				ticker.Stop()
				l.unlockWithLog(ctx)
				return
			default:
				if success := l.extend(ctx); !success {
					ticker.Stop()
					return
				}
			}
		}
	}
}
