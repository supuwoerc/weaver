package conf

import "fmt"

type Config struct {
	AppName       string              `mapstructure:"app_name"`       // 应用名称
	AppVersion    string              `mapstructure:"app_version"`    // 应用版本
	Env           string              `mapstructure:"env"`            // 环境
	System        SystemConfig        `mapstructure:"system"`         // 系统相关配置
	Email         Email               `mapstructure:"email"`          // 邮件配置
	JWT           JWTConfig           `mapstructure:"jwt"`            // jwt相关配置
	Logger        LoggerConfig        `mapstructure:"logger"`         // logger相关配置
	Cors          CorsConfig          `mapstructure:"cors"`           // cors相关配置
	Captcha       CaptchaConfig       `mapstructure:"captcha"`        // 验证码相关配置
	Account       AccountConfig       `mapstructure:"account"`        // 账户相关配置
	Consul        ConsulConfig        `mapstructure:"consul"`         // consul配置
	Redis         RedisConfig         `mapstructure:"redis"`          // redis配置
	GORM          GORMConfig          `mapstructure:"gorm"`           // gorm配置
	OSS           OSSConfig           `mapstructure:"oss"`            // oss配置
	OpenTelemetry OpenTelemetryConfig `mapstructure:"open_telemetry"` // open telemetry配置
	OTLP          OTLPConfig          `mapstructure:"otlp"`           // otlp配置
	Elasticsearch ElasticsearchConfig `mapstructure:"elasticsearch"`  // es配置
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
func (c *Config) AppInfo() string {
	return fmt.Sprintf("%s:%s:%s", c.AppName, c.AppVersion, c.Env)
}
