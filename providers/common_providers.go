package providers

import (
	goredislib "github.com/redis/go-redis/v9"
	v1 "github.com/supuwoerc/weaver/api/v1"
	"github.com/supuwoerc/weaver/initialize"
	"github.com/supuwoerc/weaver/middleware"
	"github.com/supuwoerc/weaver/pkg/captcha"
	"github.com/supuwoerc/weaver/pkg/jwt"
	"github.com/supuwoerc/weaver/pkg/logger"
	"github.com/supuwoerc/weaver/pkg/utils"
	"github.com/supuwoerc/weaver/repository"
	"github.com/supuwoerc/weaver/repository/cache"
	"github.com/supuwoerc/weaver/repository/dao"
	"github.com/supuwoerc/weaver/service"
	gormLogger "gorm.io/gorm/logger"

	"github.com/google/wire"
	"github.com/mojocn/base64Captcha"
)

var loggerProvider = wire.NewSet(
	wire.Bind(new(utils.LocksmithLogger), new(*logger.Logger)),
	wire.Bind(new(initialize.ClientLogger), new(*logger.Logger)),
	logger.NewLogger,
)

var redisLoggerProvider = wire.NewSet(
	wire.Bind(new(goredislib.Hook), new(*initialize.RedisLogger)),
	initialize.NewRedisLogger,
)

var gormLoggerProvider = wire.NewSet(
	wire.Bind(new(gormLogger.Interface), new(*initialize.GormLogger)),
	initialize.NewGormLogger,
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
	loggerProvider,
	redisLoggerProvider,
	gormLoggerProvider,
	emailProvider,
	departmentServiceProvider,
	permissionServiceProvider,
	userRepositoryProvider,
	roleRepositoryProvider,
	permissionRepositoryProvider,
	captchaRedisStoreProvider,
)
