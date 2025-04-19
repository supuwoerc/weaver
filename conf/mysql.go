package conf

import "time"

type MysqlConfig struct {
	DSN         string        `mapstructure:"dsn"`           // dsn
	MaxIdleConn int           `mapstructure:"max_idle_conn"` // max_idle_conn
	MaxOpenConn int           `mapstructure:"max_open_conn"` // max_open_conn
	MaxLifetime time.Duration `mapstructure:"max_life_time"` // max_life_time
}
