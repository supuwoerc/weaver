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
	CaptchaApiProvider,
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
	dao.NewBasicDao,
	middleware.NewAuthMiddleware,
	wire.Bind(new(middleware.AuthMiddlewareTokenRepo), new(*repository.UserRepository)),
	repository.NewUserRepository,
	wire.Bind(new(repository.UserDAO), new(*dao.UserDAO)),
	wire.Bind(new(repository.UserCache), new(*cache.UserCache)),
	dao.NewUserDAO,
	cache.NewUserCache,
	jwt.NewJwtBuilder,
	wire.Bind(new(jwt.TokenBuilderRepo), new(*repository.UserRepository)),
)

var CaptchaApiProvider = wire.NewSet(
	v1.NewCaptchaApi,
	wire.Bind(new(v1.CaptchaService), new(*service.CaptchaService)),
	service.NewCaptchaService,
)

var DepartmentApiProvider = wire.NewSet(
	v1.NewDepartmentApi,
	wire.Bind(new(v1.DepartmentService), new(*service.DepartmentService)),
	service.NewDepartmentService,
	wire.Bind(new(service.DepartmentRepository), new(*repository.DepartmentRepository)),
	wire.Bind(new(service.UserRepository), new(*repository.UserRepository)),
	repository.NewDepartmentRepository,
	wire.Bind(new(repository.DepartmentDAO), new(*dao.DepartmentDAO)),
	dao.NewDepartmentDAO,
	wire.Bind(new(repository.DepartmentCache), new(*cache.DepartmentCache)),
	cache.NewDepartmentCache,
)
