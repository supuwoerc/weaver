package utils

import (
	"context"
	"fmt"
	"gin-web/pkg/constant"
	"gin-web/pkg/redis"
	"gin-web/pkg/response"
	"github.com/go-redsync/redsync/v4"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync"
	"time"
)

type RedisLocksmith struct {
	logger      *zap.SugaredLogger
	redisClient *redis.CommonRedisClient
}

type RedisLock struct {
	*redsync.Mutex
	duration time.Duration
	logger   *zap.SugaredLogger
}

const (
	defaultMaxRetries   = 32
	defaultLockDuration = 10 * time.Second
)

var (
	redisLocksmith     *RedisLocksmith
	redisLocksmithOnce sync.Once
)

func NewRedisLocksmith(logger *zap.SugaredLogger, redisClient *redis.CommonRedisClient) *RedisLocksmith {
	redisLocksmithOnce.Do(func() {
		redisLocksmith = &RedisLocksmith{
			logger:      logger,
			redisClient: redisClient,
		}
	})
	return redisLocksmith
}

// NewLock 创建锁
func (r *RedisLocksmith) NewLock(t constant.Prefix, object ...interface{}) *RedisLock {
	var temp []string
	switch any(object).(type) {
	case []uint:
		temp = lo.Map(any(object).([]uint), func(item uint, _ int) string {
			return strconv.Itoa(int(item))
		})
	case []string:
		temp = any(object).([]string)
	default:
		// TODO: 测试default 分支
		temp = any(object).([]string)
	}
	name := fmt.Sprintf("%s:%s", t, strings.Join(temp, "_"))
	return &RedisLock{
		Mutex:    r.redisClient.Redsync.NewMutex(name, redsync.WithExpiry(defaultLockDuration), redsync.WithTries(defaultMaxRetries)),
		duration: defaultLockDuration,
		logger:   r.logger,
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
		return errors.Wrapf(err, l.Name())
	}
	if !ok {
		return fmt.Errorf("%s unlock failed", l.Name())
	}
	return nil
}

func (l *RedisLock) alarm(subject constant.Subject, lockName string, err error) {
	l.logger.Errorf("redis lock name:%s subject:%s error:%s", lockName, subject, err.Error())
	//adminEmail := viper.GetString("system.admin.email")
	// TODO:全局的告警方法
	//if e := email.NewEmailClient().SendText(adminEmail, subject, fmt.Sprintf("%s alarm: %v", lockName, err)); e != nil {
	//	s.logger.Errorf("发送邮件失败,信息:%s", e.Error())
	//}
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
