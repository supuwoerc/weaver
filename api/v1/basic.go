package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/supuwoerc/weaver/conf"
	"github.com/supuwoerc/weaver/middleware"
	"github.com/supuwoerc/weaver/pkg/logger"
)

type BasicApi struct {
	Route  *gin.RouterGroup
	Logger *logger.Logger
	Conf   *conf.Config
	Auth   *middleware.AuthMiddleware
}

func NewBasicApi(
	route *gin.RouterGroup,
	logger *logger.Logger,
	conf *conf.Config,
	auth *middleware.AuthMiddleware,
) *BasicApi {
	return &BasicApi{
		Route:  route,
		Logger: logger,
		Conf:   conf,
		Auth:   auth,
	}
}
