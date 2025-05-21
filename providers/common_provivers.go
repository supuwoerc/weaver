package providers

import (
	"github.com/google/wire"
	"github.com/supuwoerc/weaver/middleware"
	"github.com/supuwoerc/weaver/pkg/jwt"
	"github.com/supuwoerc/weaver/repository/cache"
	"github.com/supuwoerc/weaver/service"
)

var userCacheProvider = wire.NewSet(
	wire.Bind(new(middleware.AuthMiddlewareTokenRepo), new(*cache.UserCache)),
	wire.Bind(new(jwt.TokenBuilderRepo), new(*cache.UserCache)),
	wire.Bind(new(service.UserCache), new(*cache.UserCache)),
	cache.NewUserCache,
)

var CommonProvider = wire.NewSet(
	userCacheProvider,
)
