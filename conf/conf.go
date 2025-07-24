package conf

type Config struct {
	Env     string        `mapstructure:"env"`     // 环境
	System  SystemConfig  `mapstructure:"system"`  // 系统相关配置
	JWT     JWTConfig     `mapstructure:"jwt"`     // jwt相关配置
	Logger  LoggerConfig  `mapstructure:"logger"`  // logger相关配置
	Cors    CorsConfig    `mapstructure:"cors"`    // cors相关配置
	Captcha CaptchaConfig `mapstructure:"captcha"` // 验证码相关配置
	Account AccountConfig `mapstructure:"account"` // 账户相关配置
	Redis   RedisConfig   `mapstructure:"redis"`   // redis配置
	Mysql   MysqlConfig   `mapstructure:"mysql"`   // mysql配置
	OSS     string        `mapstructure:"oss"`     // oss配置
}

func (c *Config) IsProd() bool {
	return c.Env == "prod"
}
func (c *Config) IsDev() bool {
	return c.Env == "dev"
}

func (c *Config) IsTest() bool {
	return c.Env == "test"
}
