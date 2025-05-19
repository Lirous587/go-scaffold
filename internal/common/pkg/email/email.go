package email

import (
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
)

// Mailer 邮件发送接口
type Mailer interface {
	Send(to, subject string, body string) error
	SendHTML(to, subject string, htmlBody string) error
}

// GetMailerInstance 获取全局邮件客户端实例
func GetMailerInstance() Mailer {
	return instance
}

func (m *mailer) Send(to, subject, body string) error {
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

	return errors.WithStack(m.dialer.DialAndSend(msg))
}
