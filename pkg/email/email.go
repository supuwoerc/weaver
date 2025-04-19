package email

import (
	"bytes"
	"gin-web/conf"
	"gin-web/pkg/constant"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"html/template"
	"path/filepath"
)

type Client struct {
	logger *zap.SugaredLogger
	dialer *gomail.Dialer
	conf   *conf.Config
}

func NewEmailClient(logger *zap.SugaredLogger, dialer *gomail.Dialer, conf *conf.Config) *Client {
	return &Client{
		logger: logger,
		dialer: dialer,
		conf:   conf,
	}
}

func (e *Client) send(to string, subject constant.Subject, body string, c constant.MIME) error {
	if e.conf.IsProd() {
		message := gomail.NewMessage()
		message.SetHeader("From", e.dialer.Username)
		message.SetHeader("To", to)
		message.SetHeader("Subject", string(subject))
		message.SetBody(string(c), body)
		return e.dialer.DialAndSend(message)
	} else {
		e.logger.Debugw("Sending emails in non-production environments", "To", to, "Subject", string(subject), "Env", e.conf.Env)
	}
	return nil
}

func (e *Client) SendText(to string, subject constant.Subject, body string) error {
	return e.send(to, subject, body, constant.TextPlain)
}

func (e *Client) SendHTML(to string, subject constant.Subject, templatePath constant.Template, data any) error {
	dir := e.conf.System.EmailTemplateDir
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

func (e *Client) Alarm2Admin(subject constant.Subject, body string) error {
	return e.send(e.conf.System.Admin.Email, subject, body, constant.TextPlain)
}
