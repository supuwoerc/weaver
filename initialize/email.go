package initialize

import (
	"crypto/tls"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
	"strings"
)

func NewDialer() *gomail.Dialer {
	h := viper.GetString("system.email.host")
	p := viper.GetInt("system.email.port")
	u := viper.GetString("system.email.user")
	pwd := viper.GetString("system.email.password")
	if strings.TrimSpace(h) == "" {
		panic("email.host is required")
	}
	if p == 0 {
		panic("email.port is required")
	}
	if strings.TrimSpace(u) == "" {
		panic("email.user is required")
	}
	if strings.TrimSpace(pwd) == "" {
		panic("email.pwd is required")
	}
	d := gomail.NewDialer(h, p, u, pwd)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return d
}
