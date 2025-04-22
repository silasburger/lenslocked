package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type MailTrap struct {
	Token        string
	SendEndpoint string
}

func (mt MailTrap) Dial(email Email) (*http.Request, error) {
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
