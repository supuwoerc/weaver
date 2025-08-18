package providers

import (
	"github.com/google/wire"
	"github.com/supuwoerc/weaver/middleware"
	"github.com/supuwoerc/weaver/pkg/jwt"
)

var authMiddlewareProvider = wire.NewSet(
	jwt.NewJwtBuilder,
	middleware.NewAuthMiddleware,
)
var limiterMiddlewareProvider = wire.NewSet(
	middleware.NewLimiterMiddleware,
)

var MiddlewareProvider = wire.NewSet(
	authMiddlewareProvider,
	limiterMiddlewareProvider,
)
