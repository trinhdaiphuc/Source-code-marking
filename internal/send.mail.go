package internal

import (
	"net/smtp"
	"os"
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
func SendMail(from, password, to string, subject, content string) (err error) {
	// smtp server configuration.
	smtpServer := smtpServer{
		host: os.Getenv("EMAIL_SMTP_SERVER_HOST"),
		port: os.Getenv("EMAIL_SMTP_SERVER_PORT"),
	}
	// Message.
	message := []byte("To: " + to + "\r\n" + "Subject: " + subject + "\r\n" + "\r\n" + content + "\r\n")
	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpServer.host)
	toUser := []string{to}
	// Sending email.
	err = smtp.SendMail(smtpServer.Address(), auth, from, toUser, message)

	return
}
