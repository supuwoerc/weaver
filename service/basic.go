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

func NewBasicService(ctx *gin.Context) *BasicService {
	return &BasicService{
		Logger: global.Logger,
		ctx:    ctx,
	}
}
