package conf

type Email struct {
	Host     string `mapstructure:"host"`     // host
	Port     int    `mapstructure:"port"`     // port
	User     string `mapstructure:"host"`     // user
	Password string `mapstructure:"password"` // password
}
