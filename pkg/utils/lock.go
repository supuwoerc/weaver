package utils

import (
	"context"
	"errors"
	"fmt"
	"gin-web/pkg/global"
	"gin-web/pkg/response"
	"github.com/go-redsync/redsync/v4"
	"sync"
	"time"
)

type RedisLock struct {
	*redsync.Mutex
	duration time.Duration
	dog      *Watchdog
	once     sync.Once
}

const (
	defaultMaxRetries = 3
)

// NewRedisLock 创建锁
func NewRedisLock(name string, t time.Duration) *RedisLock {
	return &RedisLock{
		Mutex:    global.RedisClient.Redsync.NewMutex(name, redsync.WithExpiry(t), redsync.WithTries(defaultMaxRetries)),
		duration: t,
	}
}

// Lock 会多次尝试(defaultMaxRetries 次),如果尝试次数内还未获取到锁则返回错误
func Lock(ctx context.Context, lock *RedisLock, extend bool) error {
	err := lock.Lock()
	if err != nil {
		var e *redsync.ErrTaken
		if errors.As(err, &e) {
			return response.Busy
		}
		return fmt.Errorf("lock err: %v", err)
	}
	if extend {
		go autoExtend(ctx, lock)
	}
	return nil
}

// TryLock 获取不到锁直接返回错误
func TryLock(ctx context.Context, lock *RedisLock, extend bool) error {
	err := lock.TryLock()
	if err != nil {
		var e *redsync.ErrTaken
		if errors.As(err, &e) {
			return response.Busy
		}
		return fmt.Errorf("lock err: %v", err)
	}
	if extend {
		go autoExtend(ctx, lock)
	}
	return nil
}

func Unlock(lock *RedisLock) error {
	if lock.dog.stopChan == nil {
		return nil
	}
	ok, err := lock.Unlock()
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("%s unlock failed", lock.Name())
	}
	if lock.dog != nil {
		lock.dog.Stop()
	}
	return nil
}

func autoExtend(ctx context.Context, lock *RedisLock) {
	lock.once.Do(func() {
		lock.dog = NewWatchdog(lock.duration/2, 0, func() error {
			if ok, temp := lock.Extend(); temp != nil {
				return temp
			} else if !ok {
				return fmt.Errorf("%s extend failed", lock.Name())
			} else {
				return nil
			}
		}, func(dog *Watchdog, err error) {
			global.Logger.Errorf("%s extend lock fail: %v", lock.Name(), err)
			dog.Stop()
		})
		go func() {
			select {
			case <-ctx.Done():
				lock.dog.Stop()
			}
		}()
		lock.dog.Start()
	})
}
