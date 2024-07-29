package global

import (
	"gin-web/pkg/redis"
	"github.com/casbin/casbin/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Logger         *zap.SugaredLogger
	DB             *gorm.DB
	RedisClient    *redis.RedisClient
	Localizer      map[string]*i18n.Localizer
	LocaleErrors   map[string]map[int]error
	CasbinEnforcer *casbin.SyncedCachedEnforcer
)

const (
	CN string = "cn"
	EN string = "en"
)
