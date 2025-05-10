package notifier

import (
	"fmt"
	"github.com/pkg/errors"
	"scaffold/pkg/email"
)

type mailer struct {
	mailer email.Mailer
}

func NewMailer(m email.Mailer) Notifier {
	return &mailer{
		mailer: m,
	}
}

func (m *mailer) SendMockNotification(to string, id int) error {
	subject := "mock message with id"
	content := fmt.Sprintf(
		"<p>id：%d</p>", id,
	)

	return errors.WithStack(m.mailer.SendHTML(to, subject, content))
}
