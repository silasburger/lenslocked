package models

import (
	"fmt"

	"github.com/go-mail/mail/v2"
)

const (
	DefaultSender = "support@lenslocked.silasburger.com"
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
	ServerURL     string
}

func NewEmailService(config SMTPConfig, serverURL string) *EmailService {
	es := EmailService{
		Dialer:    mail.NewDialer(config.Host, config.Port, config.Username, config.Password),
		ServerURL: serverURL,
	}
	return &es
}

func (es EmailService) Send(email Email) error {
	msg := mail.NewMessage()
	msg.SetHeader("To", email.To)
	msg.SetHeader("Subject", email.Subject)
	es.setFrom(msg, email)
	switch {
	case email.Plaintext != "" && email.HTML != "":
		msg.SetBody("text/html", email.HTML)
		msg.AddAlternative("text/plain", email.Plaintext, mail.SetPartEncoding(mail.Base64))
	case email.Plaintext != "":
		msg.SetBody("text/plain", email.Plaintext, mail.SetPartEncoding(mail.Base64))
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
	htmlBody := fmt.Sprintf(`
	<html>
	<body>
		<p>o reset your password please visit the following URL:</p>
		<a href="http://%s">Reset Password</a>
	</body>
	</html>
	`, resetURL)
	plaintextBody := fmt.Sprintf("To reset your password please visit the following URL: %s", resetURL)
	email := Email{
		To:        to,
		Subject:   "Reset your password",
		Plaintext: plaintextBody,
		HTML:      htmlBody,
	}
	err := es.Send(email)
	if err != nil {
		return fmt.Errorf("ForgotPassword: %w", err)
	}
	return nil
}

func (es EmailService) PasswordlessSignin(to, resetURL string) error {
	htmlBody := fmt.Sprintf(`
		<html>
		<body>
			<p>Click below to sign in to your account:</p>
			<a href="http://%s">Sign in</a>
		</body>
		</html>
		`, resetURL)
	plaintextBody := fmt.Sprintf("To sign in to your account visit following URL: %s", resetURL)
	email := Email{
		To:        to,
		Subject:   "Sign in link",
		Plaintext: plaintextBody,
		HTML:      htmlBody,
	}
	err := es.Send(email)
	if err != nil {
		return fmt.Errorf("SendSignin: %w", err)
	}
	return nil
}
