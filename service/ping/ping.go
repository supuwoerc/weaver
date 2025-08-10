package ping

import (
	"context"
	"time"

	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/utils"
	"github.com/supuwoerc/weaver/service"
)

type Service struct {
	*service.BasicService
}

func NewPingService(basic *service.BasicService) *Service {
	return &Service{
		BasicService: basic,
	}
}

func (p *Service) LockPermissionField(ctx context.Context) error {
	lock := p.Locksmith.NewLock(constant.PermissionIdPrefix, "100", "200")
	if err := lock.TryLock(ctx, true); err != nil {
		return err
	}
	p.Logger.WithContext(ctx).Info("lock success")
	defer func(lock *utils.RedisLock) {
		e := lock.Unlock()
		if e != nil {
			p.Logger.WithContext(ctx).Infof("unlock fail %s", e.Error())
			return
		}
		p.Logger.WithContext(ctx).Info("unlock success")
	}(lock)
	time.Sleep(time.Second * 20)
	return nil
}
