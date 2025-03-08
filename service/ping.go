package service

import (
	"context"
	"gin-web/pkg/constant"
	pkgRedis "gin-web/pkg/redis"
	"gin-web/pkg/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
	"time"
)

type PingService struct {
	*BasicService
}

var (
	pingServiceOnce sync.Once
	pingService     *PingService
)

func NewPingService(logger *zap.SugaredLogger, db *gorm.DB, r *pkgRedis.CommonRedisClient,
	locksmith *utils.RedisLocksmith, v *viper.Viper) *PingService {
	pingServiceOnce.Do(func() {
		pingService = &PingService{
			BasicService: NewBasicService(logger, r, db, locksmith, v),
		}
	})
	return pingService
}

func (p *PingService) LockPermissionField(ctx context.Context) error {
	lock := p.locksmith.NewLock(constant.PermissionIdPrefix, 100, 200)
	if err := lock.TryLock(ctx, true); err != nil {
		return err
	}
	p.logger.Info("lock success")
	defer func(lock *utils.RedisLock) {
		e := lock.Unlock()
		if e != nil {
			p.logger.Infof("unlock fail %s", e.Error())
			return
		}
		p.logger.Info("unlock success")
	}(lock)
	time.Sleep(time.Second * 20)
	return nil
}
