package service

import (
	"gin-web/pkg/global"
	"go.uber.org/zap"
)

type BasicService struct {
	Logger *zap.SugaredLogger
}

var basicService *BasicService

func NewBasicService() *BasicService {
	if basicService == nil {
		basicService = &BasicService{
			Logger: global.Logger,
		}
	}
	return basicService
}
