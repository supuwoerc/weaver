package v1

import (
	"gin-web/pkg/global"
	"go.uber.org/zap"
)

type BasicApi struct {
	logger *zap.SugaredLogger
}

var basicApi *BasicApi

func NewBasicApi() *BasicApi {
	if basicApi == nil {
		basicApi = &BasicApi{logger: global.Logger}
	}
	return basicApi
}
