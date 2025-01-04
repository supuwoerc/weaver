package global

import (
	"gin-web/pkg/redis"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

var (
	Logger      *zap.SugaredLogger
	DB          *gorm.DB
	RedisClient *redis.RedisClient
	Localizer   map[string]*i18n.Localizer
	Dialer      *gomail.Dialer
)

const (
	CN string = "cn"
	EN string = "en"
)
