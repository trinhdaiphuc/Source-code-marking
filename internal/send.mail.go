package internal

import (
	"os"

	"github.com/mailgun/mailgun-go"
)

type smtpServer struct {
	host string
	port string
}

// Address URI to smtp server
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

// SendMail send an email from user to another users
func SendMail(from, to, subject, content string) (id string, err error) {
	mg := mailgun.NewMailgun(os.Getenv("EMAIL_MAILGUN_DOMAIN"), os.Getenv("EMAIL_MAILGUN_API_KEY"))
	m := mg.NewMessage(os.Getenv("EMAIL_USERNAME"), subject, content, to)
	_, id, err = mg.Send(m)
	return id, err
}
