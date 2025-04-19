package providers

import (
	v1 "gin-web/api/v1"
	"gin-web/middleware"
	"gin-web/repository"
	"gin-web/repository/cache"
	"gin-web/repository/dao"
	"gin-web/service"
	"github.com/google/wire"
)

// V1Provider api-provider集合
var V1Provider = wire.NewSet(
	dao.NewBasicDao,
	service.NewBasicService,
	AttachmentApiProvider,
	CaptchaApiProvider,
	DepartmentApiProvider,
	PermissionApiProvider,
	PingApiProvider,
	RoleApiProvider,
	UserApiProvider,
)

var AttachmentApiProvider = wire.NewSet(
	v1.NewAttachmentApi,
	wire.Bind(new(v1.AttachmentService), new(*service.AttachmentService)),
	service.NewAttachmentService,
	wire.Bind(new(service.AttachmentRepository), new(*repository.AttachmentRepository)),
	repository.NewAttachmentRepository,
	wire.Bind(new(repository.AttachmentDAO), new(*dao.AttachmentDAO)),
	dao.NewAttachmentDAO,
	middleware.NewAuthMiddleware,
)

var CaptchaApiProvider = wire.NewSet(
	v1.NewCaptchaApi,
	wire.Bind(new(v1.CaptchaService), new(*service.CaptchaService)),
	service.NewCaptchaService,
)

var DepartmentApiProvider = wire.NewSet(
	v1.NewDepartmentApi,
	wire.Bind(new(service.DepartmentRepository), new(*repository.DepartmentRepository)),
	repository.NewDepartmentRepository,
	wire.Bind(new(repository.DepartmentDAO), new(*dao.DepartmentDAO)),
	dao.NewDepartmentDAO,
	wire.Bind(new(repository.DepartmentCache), new(*cache.DepartmentCache)),
	cache.NewDepartmentCache,
)

var PermissionApiProvider = wire.NewSet(
	v1.NewPermissionApi,
)

var PingApiProvider = wire.NewSet(
	v1.NewPingApi,
	wire.Bind(new(v1.PingService), new(*service.PingService)),
	service.NewPingService,
)

var RoleApiProvider = wire.NewSet(
	v1.NewRoleApi,
	wire.Bind(new(v1.RoleService), new(*service.RoleService)),
	service.NewRoleService,
)

var UserApiProvider = wire.NewSet(
	v1.NewUserApi,
	wire.Bind(new(v1.UserService), new(*service.UserService)),
	service.NewUserService,
)
