package v1

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"sync"
)

type BasicApi struct {
	logger *zap.SugaredLogger
	viper  *viper.Viper
}

var (
	basicOnce sync.Once
	basicApi  *BasicApi
)

func NewBasicApi(logger *zap.SugaredLogger, v *viper.Viper) *BasicApi {
	basicOnce.Do(func() {
		basicApi = &BasicApi{
			logger: logger,
			viper:  v,
		}
	})
	return basicApi
}
