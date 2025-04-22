package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type MailConfig struct {
	Host  string
	Token string
}

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

func main() {

	htmlBody := fmt.Sprintf(`
		<html>
		<body>
			<p>Click below to sign in to your account:</p>
			<a href="https://%s">Sign in</a>
		</body>
		</html>
		`, "google.com")
	data := Body{
		To:      []Address{{Email: "ybsilas@gmail.com", Name: "Silas"}},
		From:    Address{Email: "support@lenslocked.silasburger.com", Name: "Support"},
		Subject: "TEST",
		HTML:    htmlBody,
	}

	bodyJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("malformed json")
		return
	}

	fmt.Println(bodyJSON)

	err = godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	url, err := url.Parse("https://" + os.Getenv("MAIL_HOST") + "/api/send")
	if err != nil {
		fmt.Println("invalid url")
	}

	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(bodyJSON))
	if err != nil {
		fmt.Println("failed creating request")
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Api-Token", os.Getenv("MAIL_TOKEN"))

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("request failed: %w", err)
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	fmt.Println("Status Code:", res.StatusCode)
	fmt.Println("Response Body:", string(body))
}
