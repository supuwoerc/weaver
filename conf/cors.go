package conf

type CorsConfig struct {
	OriginPrefix []string `mapstructure:"origin_prefix"` // origin前缀
}
