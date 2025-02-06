package models

import (
	"fmt"

	"github.com/go-mail/mail/v2"
)

const (
	DefaultSender = "support@lenslocked.com"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type Email struct {
	To        string
	From      string
	Subject   string
	Plaintext string
	HTML      string
}

type EmailService struct {
	DefaultSender string
	Dialer        *mail.Dialer
}

func NewEmailService(config SMTPConfig) *EmailService {
	es := EmailService{
		Dialer: mail.NewDialer(config.Host, config.Port, config.Username, config.Password),
	}
	return &es
}

func (es EmailService) Send(email Email) error {
	msg := mail.NewMessage()
	msg.SetHeader("To", email.To)
	// TODO: Set the from field to a default value if it isn't set by Email
	msg.SetHeader("Subject", email.Subject)
	es.setFrom(msg, email)
	switch {
	case email.Plaintext != "" && email.HTML != "":
		msg.SetBody("text/plain", email.Plaintext)
		msg.AddAlternative("text/html", email.HTML)
	case email.Plaintext != "":
		msg.SetBody("text/plain", email.Plaintext)
	case email.HTML != "":
		msg.SetBody("text/html", email.HTML)
	}
	err := es.Dialer.DialAndSend(msg)
	if err != nil {
		return fmt.Errorf("Send: %w", err)
	}
	return nil
}

func (es EmailService) setFrom(msg *mail.Message, email Email) {
	var from string
	switch {
	case email.From != "":
		from = email.From
	case es.DefaultSender != "":
		from = es.DefaultSender
	default:
		from = DefaultSender
	}
	msg.SetHeader("From", from)
}

func (es EmailService) ForgotPassword(to, resetURL string) error {
	email := Email{
		To:        to,
		Subject:   "Reset your password",
		Plaintext: "To reset your password please visit the following URL: " + resetURL,
		HTML:      `<p>To reset your password please visit the following URL: <a href="` + resetURL + `"> Reset Password </a></p>`,
	}
	err := es.Send(email)
	if err != nil {
		return fmt.Errorf("ForgotPassword: %w", err)
	}
	return nil
}
