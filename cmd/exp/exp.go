package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/silasburger/lenslocked/models"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	type SMTPConfig struct {
		Host     string
		Port     int
		Username string
		Password string
	}

	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}

	// config := models.SMTPConfig{
	// 	Host:     host,
	// 	Port:     port,
	// 	Username: username,
	// 	Password: password,
	// }
	// email := models.Email{
	// 	From:      "test@lenslocked.com",
	// 	To:        "jon@calhoun.io",
	// 	Subject:   "This is a test email",
	// 	Plaintext: "This is the body of the email",
	// 	HTML:      `<h1>Hello there buddy!</h1><p>This is the email</p><p>Hope you enjoy it</p>`,
	// }
	es := models.NewEmailService(models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})
	err = es.ForgotPassword("ybsilas@gmail.com", "localhost:3000/reset-pw?token=abc123")
	if err != nil {
		panic(err)
	}

	fmt.Println("email sent")
}
