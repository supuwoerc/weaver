package email

import (
	"gin-web/pkg/global"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

func isProd() bool {
	return viper.Get("env") == "prod"
}

func SendText(to, subject, body string) error {
	if isProd() {
		dialer := global.Dialer
		message := gomail.NewMessage()
		message.SetHeader("From", dialer.Username)
		message.SetHeader("To", to)
		message.SetHeader("Subject", subject)
		message.SetBody("text/plain", body)
		return dialer.DialAndSend(message)
	}
	return nil
}
