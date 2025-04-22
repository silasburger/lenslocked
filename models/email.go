package models

import (
	"fmt"
	"io"
	"net/http"
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
	EmailAPI
}

type EmailAPI interface {
	Dial(email Email) (*http.Request, error)
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
	email := Email{
		To:      to,
		Subject: "Reset your password",
		Text:    plaintextBody,
		HTML:    htmlBody,
	}
	email.From = es.setFrom(email)
	req, err := es.EmailAPI.Dial(email)
	if err != nil {
		return fmt.Errorf("ForgotPassword: %w", err)
	}
	err = es.Send(req)
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
	email := Email{
		To:      to,
		Subject: "Sign in link",
		Text:    plaintextBody,
		HTML:    htmlBody,
	}
	email.From = es.setFrom(email)
	req, err := es.EmailAPI.Dial(email)
	if err != nil {
		return fmt.Errorf("PasswordlessSignin: %w", err)
	}
	err = es.Send(req)
	if err != nil {
		return fmt.Errorf("PasswordlessSignin: %w", err)
	}
	return nil
}

func (es EmailService) Send(req *http.Request) error {
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}
	ok := res.StatusCode >= 200 && res.StatusCode < 300
	if !ok {
		body, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			return fmt.Errorf("send: %w", err)
		}
		return fmt.Errorf("send: HTTP %d: %s", res.StatusCode, string(body))
	}
	return nil
}

func (es EmailService) setFrom(email Email) string {
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
