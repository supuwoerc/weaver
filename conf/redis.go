package conf

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`     // addr
	Password string `mapstructure:"password"` // password
	DB       int    `mapstructure:"db"`       // db
}
