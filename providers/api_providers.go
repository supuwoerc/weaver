package providers

import (
	"github.com/google/wire"
	"github.com/mojocn/base64Captcha"
	v1 "github.com/supuwoerc/weaver/api/v1"
	attachment2 "github.com/supuwoerc/weaver/api/v1/attachment"
	captcha3 "github.com/supuwoerc/weaver/api/v1/captcha"
	department2 "github.com/supuwoerc/weaver/api/v1/department"
	permission2 "github.com/supuwoerc/weaver/api/v1/permission"
	ping2 "github.com/supuwoerc/weaver/api/v1/ping"
	role2 "github.com/supuwoerc/weaver/api/v1/role"
	user2 "github.com/supuwoerc/weaver/api/v1/user"
	"github.com/supuwoerc/weaver/initialize"
	"github.com/supuwoerc/weaver/pkg/captcha"
	"github.com/supuwoerc/weaver/repository/cache"
	"github.com/supuwoerc/weaver/repository/dao"
	"github.com/supuwoerc/weaver/service"
	"github.com/supuwoerc/weaver/service/attachment"
	captcha2 "github.com/supuwoerc/weaver/service/captcha"
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
	wire.Bind(new(attachment2.Service), new(*attachment.Service)),
	wire.Bind(new(attachment.DAO), new(*dao.AttachmentDAO)),
	wire.Bind(new(attachment.Storage), new(*initialize.S3CompatibleStorage)),
	initialize.NewS3CompatibleStorage,
	dao.NewAttachmentDAO,
	attachment.NewAttachmentService,
	attachment2.NewAttachmentApi,
)

var captchaApiProvider = wire.NewSet(
	wire.Bind(new(captcha3.Service), new(*captcha2.Service)),
	wire.Bind(new(base64Captcha.Store), new(*captcha.RedisStore)),
	captcha.NewRedisStore,
	captcha2.NewCaptchaService,
	captcha3.NewCaptchaApi,
)

var departmentServiceProvider = wire.NewSet(
	wire.Bind(new(DepartmentCache), new(*department.Service)),
	wire.Bind(new(department2.Service), new(*department.Service)),
	wire.Bind(new(department.DAO), new(*dao.DepartmentDAO)),
	wire.Bind(new(department.DepartmentCache), new(*cache.DepartmentCache)),
	dao.NewDepartmentDAO,
	cache.NewDepartmentCache,
	department.NewDepartmentService,
)

var departmentApiProvider = wire.NewSet(
	departmentServiceProvider,
	department2.NewDepartmentApi,
)

var roleDAOProvider = wire.NewSet(
	wire.Bind(new(permission.RoleDAO), new(*dao.RoleDAO)),
	wire.Bind(new(role.DAO), new(*dao.RoleDAO)),
	dao.NewRoleDAO,
)

var permissionDAOProvider = wire.NewSet(
	wire.Bind(new(permission.DAO), new(*dao.PermissionDAO)),
	wire.Bind(new(role.PermissionDAO), new(*dao.PermissionDAO)),
	dao.NewPermissionDAO,
)

var permissionServiceProvider = wire.NewSet(
	wire.Bind(new(PermissionCache), new(*permission.Service)),
	wire.Bind(new(permission2.Service), new(*permission.Service)),
	permissionDAOProvider,
	permission.NewPermissionService,
)

var permissionApiProvider = wire.NewSet(
	permissionServiceProvider,
	permission2.NewPermissionApi,
)

var pingApiProvider = wire.NewSet(
	wire.Bind(new(ping2.Service), new(*ping.Service)),
	ping.NewPingService,
	ping2.NewPingApi,
)

var roleApiProvider = wire.NewSet(
	wire.Bind(new(role2.Service), new(*role.Service)),
	roleDAOProvider,
	role.NewRoleService,
	role2.NewRoleApi,
)

var userApiProvider = wire.NewSet(
	wire.Bind(new(user2.Service), new(*user.Service)),
	wire.Bind(new(user.DAO), new(*dao.UserDAO)),
	dao.NewUserDAO,
	user.NewUserService,
	user2.NewUserApi,
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
