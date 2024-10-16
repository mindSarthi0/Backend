package API

import (
	"gopkg.in/gomail.v2"
	"log"
)

func sendEmailWithAttachment(to, subject, body, attachmentPath string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "your-email@example.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	m.Attach(attachmentPath)

	d := gomail.NewDialer("smtp.gmail.com", 587, "your-email@example.com", "your-password")
	err := d.DialAndSend(m)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}
	log.Printf("Email sent to %s", to)
	return nil
}
