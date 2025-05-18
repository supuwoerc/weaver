package providers

import (
	v1 "gin-web/api/v1"
	"gin-web/initialize"
	"gin-web/middleware"
	"gin-web/pkg/captcha"
	"gin-web/pkg/jwt"
	"gin-web/pkg/utils"
	"gin-web/repository"
	"gin-web/repository/cache"
	"gin-web/repository/dao"
	"gin-web/service"

	"go.uber.org/zap"

	"github.com/google/wire"
	"github.com/mojocn/base64Captcha"
)

var zapLoggerProvider = wire.NewSet(
	wire.Bind(new(utils.LocksmithLogger), new(*zap.SugaredLogger)),
	wire.Bind(new(initialize.ClientLogger), new(*zap.SugaredLogger)),
	initialize.NewZapLogger,
)

var emailProvider = wire.NewSet(
	wire.Bind(new(utils.LocksmithEmailClient), new(*initialize.EmailClient)),
	initialize.NewEmailClient,
)

var departmentServiceProvider = wire.NewSet(
	wire.Bind(new(DepartmentCache), new(*service.DepartmentService)),
	wire.Bind(new(v1.DepartmentService), new(*service.DepartmentService)),
	service.NewDepartmentService,
)

var permissionServiceProvider = wire.NewSet(
	wire.Bind(new(PermissionCache), new(*service.PermissionService)),
	wire.Bind(new(v1.PermissionService), new(*service.PermissionService)),
	service.NewPermissionService,
)

var userRepositoryProvider = wire.NewSet(
	wire.Bind(new(service.UserRepository), new(*repository.UserRepository)),
	wire.Bind(new(middleware.AuthMiddlewareTokenRepo), new(*repository.UserRepository)),
	wire.Bind(new(jwt.TokenBuilderRepo), new(*repository.UserRepository)),
	repository.NewUserRepository,
	wire.Bind(new(repository.UserDAO), new(*dao.UserDAO)),
	wire.Bind(new(repository.UserCache), new(*cache.UserCache)),
	dao.NewUserDAO,
	cache.NewUserCache,
	jwt.NewJwtBuilder,
)

var roleRepositoryProvider = wire.NewSet(
	wire.Bind(new(service.RoleRepository), new(*repository.RoleRepository)),
	repository.NewRoleRepository,
	wire.Bind(new(repository.RoleDAO), new(*dao.RoleDAO)),
	dao.NewRoleDAO,
)

var permissionRepositoryProvider = wire.NewSet(
	wire.Bind(new(service.PermissionRepository), new(*repository.PermissionRepository)),
	repository.NewPermissionRepository,
	wire.Bind(new(repository.PermissionDAO), new(*dao.PermissionDAO)),
	dao.NewPermissionDAO,
)

var captchaRedisStoreProvider = wire.NewSet(
	wire.Bind(new(base64Captcha.Store), new(*captcha.RedisStore)),
	captcha.NewRedisStore,
)

var CommonProvider = wire.NewSet(
	zapLoggerProvider,
	emailProvider,
	departmentServiceProvider,
	permissionServiceProvider,
	userRepositoryProvider,
	roleRepositoryProvider,
	permissionRepositoryProvider,
	captchaRedisStoreProvider,
)
