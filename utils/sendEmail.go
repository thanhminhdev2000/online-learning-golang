package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"online-learning-golang/models"
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

func SendResetEmail(userEmail, token string) error {
	resetLink := fmt.Sprintf("%s/reset-password/%s", os.Getenv("CLIENT_URL"), token)

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_EMAIL"))
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", "Password Reset")
	m.SetBody("text/html", fmt.Sprintf("Click <a href='%s'>here</a> to reset your password.", resetLink))

	d := gomail.NewDialer(os.Getenv("SMTP_HOST"), 587, os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_PASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendContactEmail(data models.Contact) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("Support Team <%s>", os.Getenv("SMTP_EMAIL")))
	m.SetHeader("To", "thanhminh.nguyendev@gmail.com")
	m.SetHeader("Subject", data.Title)

	body := fmt.Sprintf(`
	<p>%s</p>
	<p><strong>Liên hệ:</strong> %s</p>
	`, data.Content, data.Email)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(os.Getenv("SMTP_HOST"), 587, os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_PASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
