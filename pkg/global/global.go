package global

import (
	"gin-web/pkg/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Logger      *zap.SugaredLogger
	DB          *gorm.DB
	RedisClient *redis.RedisClient
)
