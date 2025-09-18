package conf

type OTLPConfig struct {
	Endpoint string `mapstructure:"endpoint"` // endpoint
	Insecure bool   `mapstructure:"insecure"` // 关闭TLS
}
