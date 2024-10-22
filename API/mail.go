package API

import (
	"log"

	"gopkg.in/mail.v2" // Sending email
)

func sendEmail(to string, subject string, body string, attachmentPath string) {
	m := mail.NewMessage()

	// Set email sender, receiver, subject, and body
	m.SetHeader("From", "cognify@duinvites.com") // Change this to your Zoho email
	m.SetHeader("To", to)                        // Recipient's email address
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// Attach the PDF file
	m.Attach(attachmentPath)

	// Configure Zoho SMTP settings
	d := mail.NewDialer("smtp.zoho.in", 587, "cognify@duinvites.com", "duM7zATwBkKd") //

	if err := d.DialAndSend(m); err != nil {
		log.Fatal(err)
	}

	log.Println("Email sent successfully!")
}

func SendBIG5Report(to string, attachmentPath string) {
	sendEmail(to, "Insights Unlocked: Your BIG 5 Personality Assessment Report is Ready!", "Please find attached your Extraversion report", attachmentPath)
}

func Mail() {

	pdfFile := "C:\\Users\\Rishi Raj Ganguly\\cognify-api-gateway\\output_with_background.pdf"

	recipientEmail := "blah@gmail.com"
	subject := "Your Report"
	body := "Please find attached your Extraversion report."
	sendEmail(recipientEmail, subject, body, pdfFile)
}
