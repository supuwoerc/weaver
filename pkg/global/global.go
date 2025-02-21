package global

import (
	"gin-web/pkg/redis"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

var (
	Logger      *zap.SugaredLogger
	DB          *gorm.DB
	RedisClient *redis.CommonRedisClient
	Dialer      *gomail.Dialer
	Cron        *cron.Cron
	CronLogger  cron.Logger
)
