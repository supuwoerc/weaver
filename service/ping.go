package service

import (
	"context"
	"gin-web/pkg/constant"
	"gin-web/pkg/utils"
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

func NewPingService(basic *BasicService) *PingService {
	pingServiceOnce.Do(func() {
		pingService = &PingService{
			BasicService: basic,
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
