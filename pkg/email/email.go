package email

import (
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
	"scaffold/pkg/config"
)

// Mailer 邮件发送接口
type Mailer interface {
	Send(to, subject string, body string) error
	SendHTML(to, subject string, htmlBody string) error
}

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
	CC       string // 添加抄送字段
}

type mailer struct {
	config Config
	dialer *gomail.Dialer
}

// 全局邮件客户端实例
var mailerInstance Mailer

func Init(cfg *config.EmailConfig) error {
	config := Config{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Username: cfg.Username,
		Password: cfg.Password,
		From:     cfg.From,
		FromName: cfg.FromName,
		CC:       cfg.CC,
	}

	dialer := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)

	mailerInstance = &mailer{
		config: config,
		dialer: dialer,
	}

	return nil
}

// GetMailer 获取全局邮件客户端实例
func GetMailer() Mailer {
	return mailerInstance
}

// Send 发送纯文本邮件
func (m *mailer) Send(to, subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", m.config.From, m.config.FromName)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)

	// 如果设置了抄送邮箱，则添加CC头
	if m.config.CC != "" {
		msg.SetHeader("Cc", m.config.CC)
	}

	return errors.WithStack(m.dialer.DialAndSend(msg))
}

// SendHTML 发送HTML格式邮件
func (m *mailer) SendHTML(to, subject, htmlBody string) error {
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", m.config.From, m.config.FromName)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", htmlBody)

	// 如果设置了抄送邮箱，则添加CC头
	if m.config.CC != "" {
		msg.SetHeader("Cc", m.config.CC)
	}

	return errors.WithStack(m.dialer.DialAndSend(msg))
}
