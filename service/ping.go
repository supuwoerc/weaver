package service

import (
	"context"
	"time"

	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/utils"
)

type PingService struct {
	*BasicService
}

func NewPingService(basic *BasicService) *PingService {
	return &PingService{
		BasicService: basic,
	}
}

func (p *PingService) LockPermissionField(ctx context.Context) error {
	lock := p.locksmith.NewLock(constant.PermissionIdPrefix, "100", "200")
	if err := lock.TryLock(ctx, true); err != nil {
		return err
	}
	p.logger.WithContext(ctx).Info("lock success")
	defer func(lock *utils.RedisLock) {
		e := lock.Unlock()
		if e != nil {
			p.logger.WithContext(ctx).Infof("unlock fail %s", e.Error())
			return
		}
		p.logger.WithContext(ctx).Info("unlock success")
	}(lock)
	time.Sleep(time.Second * 20)
	return nil
}
