package email

import (
	"gin-web/pkg/global"
	"gopkg.in/gomail.v2"
)

func SendText(to, subject, body string) error {
	dialer := global.Dialer
	message := gomail.NewMessage()
	message.SetHeader("From", dialer.Username)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)
	return global.Dialer.DialAndSend(message)
}
