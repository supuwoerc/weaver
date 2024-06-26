package service

import (
	"gin-web/pkg/global"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BasicService struct {
	Logger *zap.SugaredLogger
	ctx    *gin.Context
}

var basicService *BasicService

func NewBasicService(ctx *gin.Context) *BasicService {
	if basicService == nil {
		basicService = &BasicService{
			Logger: global.Logger,
			ctx:    ctx,
		}
	}
	return basicService
}
