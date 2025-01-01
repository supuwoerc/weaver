package service

import (
	"gin-web/pkg/global"
	"go.uber.org/zap"
	"sync"
)

type BasicService struct {
	logger *zap.SugaredLogger
}

var (
	basicOnce sync.Once
	basic     *BasicService
)

func NewBasicService() *BasicService {
	basicOnce.Do(func() {
		basic = &BasicService{
			logger: global.Logger,
		}
	})
	return basic
}
