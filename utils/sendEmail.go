package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

func GenerateResetToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

func SendResetEmail(userEmail string) error {
	token, err := GenerateResetToken()
	if err != nil {
		return err
	}
	resetLink := fmt.Sprintf("http://localhost:3000/request-reset-password?token=%s", token)

	m := gomail.NewMessage()
	m.SetHeader("From", "thanhminh.test01@gmail.com")
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", "Password Reset")
	m.SetBody("text/html", fmt.Sprintf("Click <a href='%s'>here</a> to reset your password.", resetLink))

	d := gomail.NewDialer(os.Getenv("SMTP_HOST"), 587, os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_PASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
