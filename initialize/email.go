package initialize

import (
	"crypto/tls"
	"gin-web/conf"
	"gopkg.in/gomail.v2"
	"strings"
)

func NewDialer(conf *conf.Config) *gomail.Dialer {
	h := conf.System.Email.Host
	p := conf.System.Email.Port
	u := conf.System.Email.User
	pwd := conf.System.Email.Password
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
