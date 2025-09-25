package email

import (
	"bytes"
	"context"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
	"html/template"
	"os"
	"strconv"
	"time"
)

type mailerConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
	CC       string
}

type mailer struct {
	dialer    *gomail.Dialer
	templates map[string]*template.Template
}

var (
	AdminEmail   string
	globalDialer *gomail.Dialer
	config       mailerConfig
)

func init() {
	if err := UpdateConfig(); err != nil {
		panic(err)
	}
	if err := validateEnv(); err != nil {
		panic(err)
	}
}

func UpdateConfig() error {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	portStr := os.Getenv("EMAIL_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return err
	}

	config = mailerConfig{
		Host:     os.Getenv("EMAIL_HOST"),
		Port:     port,
		Username: os.Getenv("EMAIL_USERNAME"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		From:     os.Getenv("EMAIL_FROM"),
		FromName: os.Getenv("EMAIL_FROM_NAME"),
		CC:       os.Getenv("EMAIL_CC"),
	}

	AdminEmail = os.Getenv("EMAIL_ADMIN")

	globalDialer = gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)

	return nil
}

func validateEnv() error {
	// 校验必填项
	if config.Host == "" ||
		config.Port == 0 ||
		config.Username == "" ||
		config.Password == "" ||
		config.From == "" ||
		AdminEmail == "" {
		return errors.New("email config: 环境变量缺失，必填项不能为空")
	}
	return nil
}

func NewMailer(templatesMap map[string]*template.Template) Mailer {
	return &mailer{
		dialer:    globalDialer,
		templates: templatesMap,
	}
}

// Mailer 邮件发送接口
type Mailer interface {
	SendPlain(to, subject, body string) error
	SendHTML(to, subject, htmlBody string) error
	SendWithTemplate(to, subject, templateName string, data ...interface{}) error
}

func (m *mailer) SendPlain(to, subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", config.From, config.FromName)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)

	// 如果设置了抄送邮箱，则添加CC头
	if config.CC != "" {
		msg.SetHeader("Cc", config.CC)
	}

	return errors.WithStack(m.dialer.DialAndSend(msg))
}

func (m *mailer) SendHTML(to, subject, htmlBody string) error {
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", config.From, config.FromName)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", htmlBody)

	// 如果设置了抄送邮箱，则添加CC头
	if config.CC != "" {
		msg.SetHeader("Cc", config.CC)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := make(chan error, 1)
	go func() {
		result <- m.dialer.DialAndSend(msg)
	}()

	select {
	case err := <-result:
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	case <-ctx.Done():
		return errors.New("发送邮件超时")
	}
}

func (m *mailer) SendWithTemplate(to, subject, templateName string, data ...interface{}) error {
	tmpl, exists := m.templates[templateName]
	if !exists {
		return errors.Errorf("模板不存在: %s", templateName)
	}

	var templateData interface{}
	// 如果有数据，使用第一个参数；否则使用 nil
	if len(data) > 0 {
		templateData = data[0]
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData); err != nil {
		return errors.Wrapf(err, "渲染模板失败: %s", templateName)
	}

	return m.SendHTML(to, subject, buf.String())
}
