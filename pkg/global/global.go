package global

import (
	"gin-web/pkg/redis"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

var (
	Logger      *zap.SugaredLogger
	DB          *gorm.DB
	RedisClient *redis.RedisClient
	Dialer      *gomail.Dialer
)
