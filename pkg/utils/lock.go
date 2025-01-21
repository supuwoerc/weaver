package utils

import (
	"context"
	"errors"
	"fmt"
	"gin-web/pkg/email"
	"gin-web/pkg/global"
	"gin-web/pkg/response"
	"github.com/go-redsync/redsync/v4"
	"github.com/spf13/viper"
	"time"
)

type RedisLock struct {
	*redsync.Mutex
	duration time.Duration
}

const (
	defaultMaxRetries = 32
)

// NewRedisLock 创建锁
func NewRedisLock(name string, t time.Duration) *RedisLock {
	return &RedisLock{
		Mutex:    global.RedisClient.Redsync.NewMutex(name, redsync.WithExpiry(t), redsync.WithTries(defaultMaxRetries)),
		duration: t,
	}
}

// Lock 会多次尝试(defaultMaxRetries 次), 如果尝试次数内还未获取到锁则返回错误
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

// Unlock 解锁
func Unlock(lock *RedisLock) error {
	ok, err := lock.Unlock()
	if err != nil {
		if errors.Is(err, redsync.ErrLockAlreadyExpired) {
			return nil
		}
		return err
	}
	if !ok {
		return fmt.Errorf("%s unlock failed", lock.Name())
	}
	return nil
}

// autoExtend 自动延长锁的过期时间
func autoExtend(ctx context.Context, lock *RedisLock) {
	ticker := time.NewTicker(lock.duration / 2)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			err := Unlock(lock)
			if err != nil {
				go func() {
					adminEmail := viper.GetString("system.admin.email")
					if e := email.SendText(adminEmail, "Unlock Fail", fmt.Sprintf("%s unlock fail: %v", lock.Name(), err)); e != nil {
						global.Logger.Errorf("发送邮件失败,信息:%s", e.Error())
					}
				}()
			}
			return
		case <-ticker.C:
			ok, err := lock.Extend()
			if err != nil {
				global.Logger.Errorf("%s extend failed: %v", lock.Name(), err)
			} else if !ok {
				global.Logger.Errorf("%s extend failed: lock couldn't be extended", lock.Name())
			}
		}
	}
}
