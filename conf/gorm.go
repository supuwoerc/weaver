package conf

import (
	"fmt"
	"time"
)

type GORMConfig struct {
	User                      string        `mapstructure:"user"`                          // user
	Password                  string        `mapstructure:"password"`                      // password
	Host                      string        `mapstructure:"host"`                          // host
	Port                      int           `mapstructure:"port"`                          // port
	Database                  string        `mapstructure:"database"`                      // database
	MaxIdleConn               int           `mapstructure:"max_idle_conn"`                 // max_idle_conn
	MaxOpenConn               int           `mapstructure:"max_open_conn"`                 // max_open_conn
	MaxLifetime               time.Duration `mapstructure:"max_life_time"`                 // max_life_time
	LogLevel                  int           `mapstructure:"log_level"`                     // gorm 日志级别
	SlowThreshold             time.Duration `mapstructure:"slow_threshold"`                // gorm 慢查询阈值
	IgnoreRecordNotFoundError bool          `mapstructure:"ignore_record_not_found_error"` // 忽略 record not found 错误
}

func (g *GORMConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", g.User, g.Password, g.Host, g.Port, g.Database)
}
