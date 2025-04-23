package models

import (
	"fmt"
)

const (
	DefaultSender = "support@lenslocked.silasburger.com"
)

type MailConfig struct {
	SendEndpoint string
	Token        string
}

type Email struct {
	To      string
	From    string
	Subject string
	Text    string
	HTML    string
}

type EmailService struct {
	DefaultSender string
	ServerURL     string
	SendEndpoint  string
	Token         string
	Emailer
}

type Emailer interface {
	DialAndSend(email *Email) error
}

func NewEmailService(emailer Emailer, serverURL string) *EmailService {
	es := EmailService{
		Emailer:   emailer,
		ServerURL: serverURL,
	}
	return &es
}
func (es EmailService) ForgotPassword(to, resetURL string) error {
	htmlBody := fmt.Sprintf(`
	<html>
	<body>
		<p>o reset your password please visit the following URL:</p>
		<a href="%s">Reset Password</a>
	</body>
	</html>
	`, resetURL)
	plaintextBody := fmt.Sprintf("To reset your password please visit the following URL: %s", resetURL)
	email := &Email{
		To:      to,
		Subject: "Reset your password",
		Text:    plaintextBody,
		HTML:    htmlBody,
	}
	email.From = es.setFrom(email)
	err := es.Send(email)
	if err != nil {
		return fmt.Errorf("ForgotPassword: %w", err)
	}
	return nil
}

func (es EmailService) PasswordlessSignin(to, signinURL string) error {
	htmlBody := fmt.Sprintf(`
		<html>
		<body>
			<p>Click below to sign in to your account:</p>
			<a href="%s">Sign in</a>
		</body>
		</html>
		`, signinURL)
	plaintextBody := fmt.Sprintf("To sign in to your account visit following URL: %s", signinURL)
	email := &Email{
		To:      to,
		Subject: "Sign in link",
		Text:    plaintextBody,
		HTML:    htmlBody,
	}
	email.From = es.setFrom(email)
	err := es.Send(email)
	if err != nil {
		return fmt.Errorf("PasswordlessSignin: %w", err)
	}
	return nil
}

func (es EmailService) Send(email *Email) error {
	err := es.Emailer.DialAndSend(email)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}
	return nil
}

func (es EmailService) setFrom(email *Email) string {
	var from string
	switch {
	case email.From != "":
		from = email.From
	case es.DefaultSender != "":
		from = es.DefaultSender
	default:
		from = DefaultSender
	}
	return from
}
