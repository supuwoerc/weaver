package v1

import (
	"gin-web/pkg/global"
	"go.uber.org/zap"
	"sync"
)

type BasicApi struct {
	logger *zap.SugaredLogger
}

var (
	basicOnce sync.Once
	basicApi  *BasicApi
)

func NewBasicApi() *BasicApi {
	basicOnce.Do(func() {
		basicApi = &BasicApi{
			logger: global.Logger,
		}
	})
	return basicApi
}
