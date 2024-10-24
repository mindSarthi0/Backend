package API

import (
	"fmt"
	"log"

	"gopkg.in/mail.v2" // Sending email
)

func sendEmail(to string, subject string, body string, attachmentPath string) {
	m := mail.NewMessage()

	// Set email sender, receiver, subject, and body
	m.SetHeader("From", "cognify@duinvites.com") // Change this to your Zoho email
	m.SetHeader("To", to)                        // Recipient's email address
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Attach the PDF file
	m.Attach(attachmentPath)

	// Configure Zoho SMTP settings
	d := mail.NewDialer("smtp.zoho.in", 587, "cognify@duinvites.com", "duM7zATwBkKd") //

	if err := d.DialAndSend(m); err != nil {
		log.Fatal(err)
	}

	log.Println("Email sent successfully!")
}

func SendBIG5Report(to string, name string, attachmentPath string) {

	htmlBody := fmt.Sprintf(`
      <p style="color: black;">Dear %s,</p>
      <p style="color: black;">We hope this email finds you well. Thank you for completing the Big 5 Personality Test! We have attached your personalized report, which provides detailed insights into your personality.</p>
      <p style="color: black;"><strong>What’s inside your report:</strong></p>
      <ul style="color: black;">
        <li>A breakdown of your scores for each of the five personality traits</li>
        <li>Personalized insights based on your responses</li>
        <li>Practical tips for personal growth and development</li>
      </ul>
      <p style="color: black;">By understanding your personality, you can gain valuable insights that can enhance your personal and professional life.</p>
      <p style="color: black;">Thank you for choosing <strong>COGNIFY</strong>. We wish you the best in your journey of self-discovery!</p>
      <p style="color: black;">Best regards,<br>Nitish</p>
    `, name)

	sendEmail(to, "Insights Unlocked: Your BIG 5 Personality Assessment Report is Ready!", htmlBody, attachmentPath)
}

func Mail() {

	pdfFile := "report.pdf"

	recipientEmail := "blah@gmail.com"
	subject := "Your Report"
	body := "Hey, here's your report!"
	sendEmail(recipientEmail, subject, body, pdfFile)
}
