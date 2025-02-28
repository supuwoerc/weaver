package utils

import (
	"context"
	"fmt"
	"gin-web/pkg/constant"
	"gin-web/pkg/email"
	"gin-web/pkg/global"
	"gin-web/pkg/response"
	"github.com/go-redsync/redsync/v4"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

type RedisLock struct {
	*redsync.Mutex
	duration time.Duration
	logger   *zap.SugaredLogger
}

const (
	defaultMaxRetries   = 32
	defaultLockDuration = 10 * time.Second
)

// NewLock 创建锁
func NewLock[T uint | string](t constant.Prefix, object ...T) *RedisLock {
	var temp []string
	switch any(object).(type) {
	case []uint:
		temp = lo.Map(any(object).([]uint), func(item uint, _ int) string {
			return strconv.Itoa(int(item))
		})
	case []string:
		temp = any(object).([]string)
	}
	name := fmt.Sprintf("%s:%s", t, strings.Join(temp, "_"))
	return &RedisLock{
		Mutex:    global.RedisClient.Redsync.NewMutex(name, redsync.WithExpiry(defaultLockDuration), redsync.WithTries(defaultMaxRetries)),
		duration: defaultLockDuration,
		logger:   global.Logger,
	}
}

// Lock 会多次尝试(defaultMaxRetries 次), 如果尝试次数内还未获取到锁则返回错误
func Lock(ctx context.Context, lock *RedisLock) error {
	err := lock.Lock()
	if err != nil {
		var e *redsync.ErrTaken
		if errors.As(err, &e) {
			return response.Busy
		}
		return fmt.Errorf("redis lock err: %v", err)
	}
	go autoExtend(ctx, lock)
	return nil
}

// TryLock 获取不到锁直接返回错误
func TryLock(ctx context.Context, lock *RedisLock) error {
	err := lock.TryLock()
	if err != nil {
		var e *redsync.ErrTaken
		if errors.As(err, &e) {
			return response.Busy
		}
		return fmt.Errorf("redis lock err: %v", err)
	}
	go autoExtend(ctx, lock)
	return nil
}

// Unlock 解锁
func Unlock(lock *RedisLock) error {
	ok, err := lock.Unlock()
	if err != nil {
		if errors.Is(err, redsync.ErrLockAlreadyExpired) {
			return nil
		}
		return errors.Wrapf(err, lock.Name())
	}
	if !ok {
		return fmt.Errorf("%s unlock failed", lock.Name())
	}
	return nil
}

func (s *RedisLock) alarm(subject constant.Subject, lockName string, err error) {
	s.logger.Errorf("redis lock name:%s subject:%s error:%s", lockName, subject, err.Error())
	adminEmail := viper.GetString("system.admin.email")
	if e := email.NewEmailClient().SendText(adminEmail, subject, fmt.Sprintf("%s alarm: %v", lockName, err)); e != nil {
		s.logger.Errorf("发送邮件失败,信息:%s", e.Error())
	}
}

func unlockWithAlarm(lock *RedisLock) {
	err := Unlock(lock)
	if err != nil {
		go lock.alarm(constant.UnlockFail, lock.Name(), err)
	}
}

func extendWithAlarm(lock *RedisLock) bool {
	ok, err := lock.Extend()
	if err != nil {
		go lock.alarm(constant.ExtendErr, lock.Name(), err)
		return false
	} else if !ok {
		go lock.alarm(constant.ExtendFail, lock.Name(), errors.New("lock couldn't be extended"))
		return false
	}
	return true
}

// autoExtend 自动延长锁的过期时间
func autoExtend(ctx context.Context, lock *RedisLock) {
	ticker := time.NewTicker(lock.duration / 2)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			unlockWithAlarm(lock)
			return
		case <-ticker.C:
			select {
			case <-ctx.Done():
				ticker.Stop()
				unlockWithAlarm(lock)
				return
			default:
				if success := extendWithAlarm(lock); !success {
					ticker.Stop()
					return
				}
			}
		}
	}
}
