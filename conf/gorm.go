package conf

import "time"

type GORMConfig struct {
	DSN                       string        `mapstructure:"dsn"`                           // dsn
	MaxIdleConn               int           `mapstructure:"max_idle_conn"`                 // max_idle_conn
	MaxOpenConn               int           `mapstructure:"max_open_conn"`                 // max_open_conn
	MaxLifetime               time.Duration `mapstructure:"max_life_time"`                 // max_life_time
	LogLevel                  int           `mapstructure:"log_level"`                     // gorm日志级别
	SlowThreshold             time.Duration `mapstructure:"slow_threshold"`                // gorm慢查询阈值
	IgnoreRecordNotFoundError bool          `mapstructure:"ignore_record_not_found_error"` // gorm忽略record not found错误
}
