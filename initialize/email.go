package initialize

import (
	"bytes"
	"context"
	"crypto/tls"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/supuwoerc/weaver/conf"
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/logger"

	"gopkg.in/gomail.v2"
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

type ClientLogger interface {
	logger.LogCtxInterface
	Debugw(msg string, keysAndValues ...interface{})
}

type EmailClient struct {
	logger ClientLogger
	dialer *gomail.Dialer
	conf   *conf.Config
}

func NewEmailClient(logger ClientLogger, dialer *gomail.Dialer, conf *conf.Config) *EmailClient {
	return &EmailClient{
		logger: logger,
		dialer: dialer,
		conf:   conf,
	}
}

func (e *EmailClient) send(ctx context.Context, to string, subject constant.Subject, body string, c constant.MIME) error {
	if e.conf.IsProd() {
		message := gomail.NewMessage()
		message.SetHeader("From", e.dialer.Username)
		message.SetHeader("To", to)
		message.SetHeader("Subject", string(subject))
		message.SetBody(string(c), body)
		return e.dialer.DialAndSend(message)
	} else {
		e.logger.WithContext(ctx).Debugw("Sending emails in non-production environments",
			"To", to, "Subject", string(subject), "Env", e.conf.Env)
	}
	return nil
}

func (e *EmailClient) SendText(ctx context.Context, to string, subject constant.Subject, body string) error {
	return e.send(ctx, to, subject, body, constant.TextPlain)
}

func (e *EmailClient) SendHTML(ctx context.Context, to string, subject constant.Subject, templatePath constant.Template, data any) error {
	dir := e.conf.System.EmailTemplateDir
	tmpl, err := template.ParseFiles(filepath.Join(dir, string(filepath.Separator), string(templatePath)))
	if err != nil {
		return err
	}
	buffer := bytes.Buffer{}
	if err = tmpl.Execute(&buffer, data); err != nil {
		return err
	}
	return e.send(ctx, to, subject, buffer.String(), constant.TextHTML)
}

func (e *EmailClient) Alarm2Admin(ctx context.Context, subject constant.Subject, body string) error {
	return e.send(ctx, e.conf.System.Admin.Email, subject, body, constant.TextPlain)
}
