package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type MailTrap struct {
	Token        string
	SendEndpoint string
}

func (mt MailTrap) Dial(email *Email) (*http.Request, error) {
	type Address struct {
		Email string
		Name  string
	}
	type Body struct {
		To      []Address
		From    Address
		Subject string
		Text    string
		HTML    string
	}
	var body Body
	body.Subject = email.Subject
	body.Text = email.Text
	body.HTML = email.HTML
	to := Address{Email: email.To, Name: ""}
	body.To = []Address{to}
	body.From = Address{Email: email.From, Name: ""}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}
	url, err := url.Parse(mt.SendEndpoint)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}
	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(bodyJSON))
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Api-Token", mt.Token)
	return req, nil
}

func (mt MailTrap) Send(emailRequest *http.Request) error {
	client := http.DefaultClient
	res, err := client.Do(emailRequest)
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

func (mt MailTrap) DialAndSend(email *Email) error {
	emailRequest, err := mt.Dial(email)
	if err != nil {
		return fmt.Errorf("dialandsend: %w", err)
	}
	err = mt.Send(emailRequest)
	if err != nil {
		return fmt.Errorf("dialandsend: %w", err)
	}
	return nil
}
