package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/supuwoerc/weaver/conf"
	"github.com/supuwoerc/weaver/middleware"
	"github.com/supuwoerc/weaver/pkg/logger"
)

type BasicApi struct {
	route  *gin.RouterGroup
	logger *logger.Logger
	conf   *conf.Config
	auth   *middleware.AuthMiddleware
}

func NewBasicApi(
	route *gin.RouterGroup,
	logger *logger.Logger,
	conf *conf.Config,
	auth *middleware.AuthMiddleware,
) *BasicApi {
	return &BasicApi{
		route:  route,
		logger: logger,
		conf:   conf,
		auth:   auth,
	}
}
