package api

import (
	v1 "gin-web/api/v1"
	"gin-web/middleware"
	"gin-web/pkg/jwt"
	"gin-web/repository"
	"gin-web/repository/cache"
	"gin-web/repository/dao"
	"gin-web/service"
	"github.com/google/wire"
)

// V1Provider api-provider集合
var V1Provider = wire.NewSet(
	AttachmentApiProvider,
)

var authMiddlewareProvider = wire.NewSet(
	middleware.NewAuthMiddleware,
	jwt.NewJwtBuilder,
	wire.Bind(new(jwt.TokenBuilderRepo), new(*repository.UserRepository)),
	wire.Bind(new(middleware.AuthMiddlewareTokenRepo), new(*repository.UserRepository)),
	repository.NewUserRepository,
	wire.Bind(new(repository.UserDAO), new(*dao.UserDAO)),
	dao.NewBasicDao,
	dao.NewUserDAO,
	wire.Bind(new(repository.UserCache), new(*cache.UserCache)),
	cache.NewUserCache,
)

var AttachmentApiProvider = wire.NewSet(
	v1.NewAttachmentApi,
	wire.Bind(new(v1.AttachmentService), new(*service.AttachmentService)),
	service.NewAttachmentService,
	service.NewBasicService,
	wire.Bind(new(service.AttachmentRepository), new(*repository.AttachmentRepository)),
	repository.NewAttachmentRepository,
	wire.Bind(new(repository.AttachmentDAO), new(*dao.AttachmentDAO)),
	dao.NewAttachmentDAO,
	authMiddlewareProvider,
)
