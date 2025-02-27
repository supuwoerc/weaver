package email

import (
	"bytes"
	"gin-web/pkg/constant"
	"gin-web/pkg/global"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"html/template"
	"path/filepath"
)

type EmailClient struct {
	logger *zap.SugaredLogger
}

func NewEmailClient() *EmailClient {
	return &EmailClient{
		logger: global.Logger,
	}
}

func (e *EmailClient) isProd() bool {
	return viper.GetString("env") == "prod"
}

func (e *EmailClient) send(to string, subject constant.Subject, body string, c constant.MIME) error {
	if e.isProd() {
		dialer := global.Dialer
		message := gomail.NewMessage()
		message.SetHeader("From", dialer.Username)
		message.SetHeader("To", to)
		message.SetHeader("Subject", string(subject))
		message.SetBody(string(c), body)
		return dialer.DialAndSend(message)
	} else {
		e.logger.Debugw("Sending emails in non-production environments", "To", to, "Subject", string(subject))
	}
	return nil
}

func (e *EmailClient) SendText(to string, subject constant.Subject, body string) error {
	return e.send(to, subject, body, constant.TextPlain)
}

func (e *EmailClient) SendHTML(to string, subject constant.Subject, templatePath constant.Template, data any) error {
	dir := viper.GetString("system.emailTemplateDir")
	tmpl, err := template.ParseFiles(filepath.Join(dir, string(filepath.Separator), string(templatePath)))
	if err != nil {
		return err
	}
	buffer := bytes.Buffer{}
	if err = tmpl.Execute(&buffer, data); err != nil {
		return err
	}
	return e.send(to, subject, buffer.String(), constant.TextHTML)
}
