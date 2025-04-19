package conf

import "time"

type AccountConfig struct {
	Expiration time.Duration `mapstructure:"expiration"` // 过期时长(秒)
}
