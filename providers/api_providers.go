package providers

import (
	"github.com/google/wire"
	"github.com/mojocn/base64Captcha"
	v1 "github.com/supuwoerc/weaver/api/v1"
	attachmentApi "github.com/supuwoerc/weaver/api/v1/attachment"
	captchaApi "github.com/supuwoerc/weaver/api/v1/captcha"
	departmentApi "github.com/supuwoerc/weaver/api/v1/department"
	permissionApi "github.com/supuwoerc/weaver/api/v1/permission"
	pingApi "github.com/supuwoerc/weaver/api/v1/ping"
	roleApi "github.com/supuwoerc/weaver/api/v1/role"
	userApi "github.com/supuwoerc/weaver/api/v1/user"
	"github.com/supuwoerc/weaver/initialize"
	"github.com/supuwoerc/weaver/middleware"
	"github.com/supuwoerc/weaver/pkg/captcha"
	"github.com/supuwoerc/weaver/repository/cache"
	"github.com/supuwoerc/weaver/repository/dao"
	"github.com/supuwoerc/weaver/service"
	"github.com/supuwoerc/weaver/service/attachment"
	captchaService "github.com/supuwoerc/weaver/service/captcha"
	"github.com/supuwoerc/weaver/service/department"
	"github.com/supuwoerc/weaver/service/permission"
	"github.com/supuwoerc/weaver/service/ping"
	"github.com/supuwoerc/weaver/service/role"
	"github.com/supuwoerc/weaver/service/user"
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
	wire.Bind(new(attachmentApi.Service), new(*attachment.Service)),
	wire.Bind(new(attachment.DAO), new(*dao.AttachmentDAO)),
	wire.Bind(new(attachment.Storage), new(*initialize.S3CompatibleStorage)),
	initialize.NewS3CompatibleStorage,
	dao.NewAttachmentDAO,
	attachment.NewAttachmentService,
	attachmentApi.NewAttachmentApi,
)

var captchaApiProvider = wire.NewSet(
	wire.Bind(new(captchaApi.Service), new(*captchaService.Service)),
	wire.Bind(new(base64Captcha.Store), new(*captcha.RedisStore)),
	captcha.NewRedisStore,
	captchaService.NewCaptchaService,
	captchaApi.NewCaptchaApi,
)

var departmentServiceProvider = wire.NewSet(
	wire.Bind(new(DepartmentCache), new(*department.Service)),
	wire.Bind(new(departmentApi.Service), new(*department.Service)),
	wire.Bind(new(department.DAO), new(*dao.DepartmentDAO)),
	wire.Bind(new(department.Cache), new(*cache.DepartmentCache)),
	dao.NewDepartmentDAO,
	cache.NewDepartmentCache,
	department.NewDepartmentService,
)

var departmentApiProvider = wire.NewSet(
	departmentServiceProvider,
	departmentApi.NewDepartmentApi,
)

var roleDAOProvider = wire.NewSet(
	wire.Bind(new(permission.RoleDAO), new(*dao.RoleDAO)),
	wire.Bind(new(role.DAO), new(*dao.RoleDAO)),
	dao.NewRoleDAO,
)

var permissionDAOProvider = wire.NewSet(
	wire.Bind(new(permission.DAO), new(*dao.PermissionDAO)),
	wire.Bind(new(role.PermissionDAO), new(*dao.PermissionDAO)),
	wire.Bind(new(middleware.AuthMiddlewarePermissionRepo), new(*dao.PermissionDAO)),
	dao.NewPermissionDAO,
)

var permissionServiceProvider = wire.NewSet(
	wire.Bind(new(PermissionCache), new(*permission.Service)),
	wire.Bind(new(permissionApi.Service), new(*permission.Service)),
	permissionDAOProvider,
	permission.NewPermissionService,
)

var permissionApiProvider = wire.NewSet(
	permissionServiceProvider,
	permissionApi.NewPermissionApi,
)

var pingApiProvider = wire.NewSet(
	wire.Bind(new(pingApi.Service), new(*ping.Service)),
	ping.NewPingService,
	pingApi.NewPingApi,
)

var roleApiProvider = wire.NewSet(
	wire.Bind(new(roleApi.Service), new(*role.Service)),
	roleDAOProvider,
	role.NewRoleService,
	roleApi.NewRoleApi,
)

var userApiProvider = wire.NewSet(
	wire.Bind(new(userApi.Service), new(*user.Service)),
	wire.Bind(new(user.DAO), new(*dao.UserDAO)),
	dao.NewUserDAO,
	user.NewUserService,
	userApi.NewUserApi,
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
