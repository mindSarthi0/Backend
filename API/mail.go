package API

import (
	"gopkg.in/mail.v2" // Sending email
	"log"
)

func sendEmail(to string, subject string, body string, attachmentPath string) {
	m := mail.NewMessage()

	// Set email sender, receiver, subject, and body
	m.SetHeader("From", "care@duinvites.com") // Change this to your Zoho email
	m.SetHeader("To", to)                     // Recipient's email address
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// Attach the PDF file
	m.Attach(attachmentPath)

	// Configure Zoho SMTP settings
	d := mail.NewDialer("smtp.zoho.com", 587, "care@duinvites.com", "Qwqw12#") //4thCpH2220XW

	if err := d.DialAndSend(m); err != nil {
		log.Fatal(err)
	}

	log.Println("Email sent successfully!")
}

func Mail() {
	pdfFile := generatePDF()

	recipientEmail := "nitishprakashb@gmail.com"
	subject := "Your Report"
	body := "Please find attached your Extraversion report."
	sendEmail(recipientEmail, subject, body, pdfFile)
}
