package providers

import (
	"github.com/google/wire"
	"github.com/mojocn/base64Captcha"
	v1 "github.com/supuwoerc/weaver/api/v1"
	"github.com/supuwoerc/weaver/pkg/captcha"
	"github.com/supuwoerc/weaver/repository/cache"
	"github.com/supuwoerc/weaver/repository/dao"
	"github.com/supuwoerc/weaver/service"
)

var basicApiProvider = wire.NewSet(
	v1.NewBasicApi,
)

var basicDAOProvider = wire.NewSet(
	dao.NewBasicDao,
)

var basicServiceProvider = wire.NewSet(
	service.NewBasicService,
)

var attachmentApiProvider = wire.NewSet(
	wire.Bind(new(v1.AttachmentService), new(*service.AttachmentService)),
	wire.Bind(new(service.AttachmentDAO), new(*dao.AttachmentDAO)),
	dao.NewAttachmentDAO,
	service.NewAttachmentService,
	v1.NewAttachmentApi,
)

var captchaApiProvider = wire.NewSet(
	wire.Bind(new(v1.CaptchaService), new(*service.CaptchaService)),
	wire.Bind(new(base64Captcha.Store), new(*captcha.RedisStore)),
	captcha.NewRedisStore,
	service.NewCaptchaService,
	v1.NewCaptchaApi,
)

var departmentServiceProvider = wire.NewSet(
	wire.Bind(new(DepartmentCache), new(*service.DepartmentService)),
	wire.Bind(new(v1.DepartmentService), new(*service.DepartmentService)),
	wire.Bind(new(service.DepartmentDAO), new(*dao.DepartmentDAO)),
	wire.Bind(new(service.DepartmentCache), new(*cache.DepartmentCache)),
	dao.NewDepartmentDAO,
	cache.NewDepartmentCache,
	service.NewDepartmentService,
)

var departmentApiProvider = wire.NewSet(
	departmentServiceProvider,
	v1.NewDepartmentApi,
)

var permissionServiceProvider = wire.NewSet(
	wire.Bind(new(PermissionCache), new(*service.PermissionService)),
	wire.Bind(new(v1.PermissionService), new(*service.PermissionService)),
	wire.Bind(new(service.PermissionDAO), new(*dao.PermissionDAO)),
	dao.NewPermissionDAO,
	service.NewPermissionService,
)

var permissionApiProvider = wire.NewSet(
	permissionServiceProvider,
	v1.NewPermissionApi,
)

var pingApiProvider = wire.NewSet(
	wire.Bind(new(v1.PingService), new(*service.PingService)),
	service.NewPingService,
	v1.NewPingApi,
)

var roleApiProvider = wire.NewSet(
	wire.Bind(new(v1.RoleService), new(*service.RoleService)),
	wire.Bind(new(service.RoleDAO), new(*dao.RoleDAO)),
	dao.NewRoleDAO,
	service.NewRoleService,
	v1.NewRoleApi,
)

var userApiProvider = wire.NewSet(
	wire.Bind(new(v1.UserService), new(*service.UserService)),
	wire.Bind(new(service.UserDAO), new(*dao.UserDAO)),
	dao.NewUserDAO,
	service.NewUserService,
	v1.NewUserApi,
)

var ApiProvider = wire.NewSet(
	basicApiProvider,
	basicServiceProvider,
	basicDAOProvider,
	attachmentApiProvider,
	captchaApiProvider,
	departmentApiProvider,
	permissionApiProvider,
	pingApiProvider,
	roleApiProvider,
	userApiProvider,
)
