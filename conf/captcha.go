package conf

import "time"

type CaptchaConfig struct {
	Expiration time.Duration `mapstructure:"expiration"` // 过期时长(秒)
}
