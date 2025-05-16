package utils

import (
	"context"
	"fmt"
	"gin-web/pkg/constant"
	"gin-web/pkg/email"
	"gin-web/pkg/redis"
	"gin-web/pkg/response"
	"github.com/go-redsync/redsync/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strings"
	"time"
)

type RedisLocksmith struct {
	logger      *zap.SugaredLogger
	redisClient *redis.CommonRedisClient
	emailClient *email.Client
}

type RedisLock struct {
	*redsync.Mutex
	duration    time.Duration
	logger      *zap.SugaredLogger
	emailClient *email.Client
}

const (
	defaultMaxRetries   = 32
	defaultLockDuration = 10 * time.Second
)

func NewRedisLocksmith(logger *zap.SugaredLogger, redisClient *redis.CommonRedisClient, emailClient *email.Client) *RedisLocksmith {
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
		Mutex:       r.redisClient.Redsync.NewMutex(name, redsync.WithExpiry(defaultLockDuration), redsync.WithTries(defaultMaxRetries)),
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

func (l *RedisLock) alarm(subject constant.Subject, lockName string, err error) {
	l.logger.Errorf("redis lock name:%s subject:%s error:%s", lockName, subject, err.Error())
	if err = l.emailClient.Alarm2Admin(subject, fmt.Sprintf("%s alarm: %v", lockName, err.Error())); err != nil {
		l.logger.Errorf("redis alarm err: %v", err.Error())
	}
}

func (l *RedisLock) unlockWithAlarm() {
	err := l.Unlock()
	if err != nil {
		go l.alarm(constant.UnlockFail, l.Name(), err)
	}
}

func (l *RedisLock) extendWithAlarm() bool {
	ok, err := l.Extend()
	if err != nil {
		go l.alarm(constant.ExtendErr, l.Name(), err)
		return false
	} else if !ok {
		go l.alarm(constant.ExtendFail, l.Name(), errors.New("lock couldn't be extended"))
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
			l.unlockWithAlarm()
			return
		case <-ticker.C:
			select {
			case <-ctx.Done():
				ticker.Stop()
				l.unlockWithAlarm()
				return
			default:
				if success := l.extendWithAlarm(); !success {
					ticker.Stop()
					return
				}
			}
		}
	}
}
