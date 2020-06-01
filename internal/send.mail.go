package internal

import (
	"net/smtp"
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
	smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}
	// Message.
	message := []byte("To: bigphuc1@gmail.com\r\n" + "Subject: " + subject + "\r\n" + "\r\n" + content + "\r\n")
	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpServer.host)
	toUser := []string{to}
	// Sending email.
	err = smtp.SendMail(smtpServer.Address(), auth, from, toUser, message)

	return
}
