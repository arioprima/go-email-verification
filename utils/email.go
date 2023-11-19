package utils

import (
	"crypto/tls"
	"fmt"
	"golang_email_verification/initializers"
	"golang_email_verification/models"
	"gopkg.in/gomail.v2"
	"log"
	"math/rand"
)

// 👇 Email template parser

func GenerateOTP() string {
	number := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = number[rand.Intn(len(number))]
	}
	return string(b)
}

func SendEmail(user *models.User, otp string) {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("could not load config", err)
	}

	// Sender data.
	from := config.EmailFrom
	smtpPass := config.SMTPPass
	smtpUser := config.SMTPUser
	to := user.Email
	smtpHost := config.SMTPHost
	smtpPort := config.SMTPPort

	// Create the email body
	body := fmt.Sprintf("Your OTP is: %s", otp)

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "OTP Verification")
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Could not send email: ", err)
	}

	fmt.Println("Email sent successfully.")
}
