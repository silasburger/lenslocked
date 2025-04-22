package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
}

func (es EmailService) Send(email Email) error {
	type Address struct {
		Email string
		Name  string
	}

	type Data struct {
		To      []Address
		From    Address
		Subject string
		Text    string
		HTML    string
	}
	var body Data
	body.Subject = email.Subject
	body.Text = email.Text
	body.HTML = email.HTML
	to := Address{Email: email.To, Name: ""}
	body.To = []Address{to}
	body.From = Address{Email: es.setFrom(email), Name: ""}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("Send: %w", err)
	}
	url, err := url.Parse(es.SendEndpoint)
	if err != nil {
		return fmt.Errorf("Send: %w", err)
	}
	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(bodyJSON))
	if err != nil {
		return fmt.Errorf("Send: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Api-Token", es.Token)
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Send: %w", err)
	}
	ok := res.StatusCode >= 200 && res.StatusCode < 300
	if !ok {
		body, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			return fmt.Errorf("send: %w", err)
		}
		return fmt.Errorf("Send: HTTP %d: %s", res.StatusCode, string(body))
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
	email := Email{
		To:      to,
		Subject: "Sign in link",
		Text:    plaintextBody,
		HTML:    htmlBody,
	}
	err := es.Send(email)
	if err != nil {
		return fmt.Errorf("SendSignin: %w", err)
	}
	return nil
}
