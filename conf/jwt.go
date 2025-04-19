package conf

import "time"

type JWTConfig struct {
	Expires             time.Duration `mapstructure:"expires"`               // token过期时长(分钟)
	RefreshTokenExpires time.Duration `mapstructure:"refresh_token_expires"` // refresh_token的过期时长(分钟)
	Secret              string        `mapstructure:"secret"`                // 密钥
	Issuer              string        `mapstructure:"issuer"`                // issuer
	TokenKey            string        `mapstructure:"token_key"`             // 客户端token对应的header-key
	RefreshTokenKey     string        `mapstructure:"refresh_token_key"`     // 客户端token对应的header-key
	TokenPrefix         string        `mapstructure:"token_prefix"`          // token前缀
}
