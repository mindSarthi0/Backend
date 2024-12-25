package API

import (
	"fmt"
	"log"
	"net"
	"time"

	"gopkg.in/mail.v2" // Sending email
)

func sendEmail(to string, subject string, body string, attachmentPath string) error {
	m := mail.NewMessage()

	// Set email sender, receiver, subject, and body
	m.SetHeader("From", "report@mindsarthi.com") // Change this to your Zoho email
	m.SetHeader("To", to)                        // Recipient's email address
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Attach the PDF file
	m.Attach(attachmentPath)

	// Configure Zoho SMTP settings
	d := mail.NewDialer("smtp.zoho.in", 587, "report@mindsarthi.com", "rXNfpN8Aa1q8") //

	err := sendEmailWithRetry(m, 5, d)

	return err
}

func sendEmailWithRetry(m *mail.Message, retries int, d *mail.Dialer) error {
	for attempt := 1; attempt <= retries; attempt++ {
		err := d.DialAndSend(m)
		if err == nil {
			log.Printf("Email sent successfully on attempt %d", attempt)
			return nil
		}

		log.Printf("Attempt %d: Error sending email: %v", attempt, err)
		if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
			log.Println("Temporary error, retrying...")
			time.Sleep(4 * time.Second)
			continue
		}

		logDetailedError(err)
		if attempt < retries {
			time.Sleep(2 * time.Second) // Wait before retrying
		}
	}
	log.Println("Failed to send email after retries")

	return fmt.Errorf("Failed to send email after retries")
}

func logDetailedError(err error) {
	log.Printf("Error sending email: %v", err)

	// Check if it's a network error to get more specific information
	if netErr, ok := err.(net.Error); ok {
		log.Printf("Is temporary: %v", netErr.Temporary())
		log.Printf("Is timeout: %v", netErr.Timeout())
	}
}

func SendBIG5Report(to string, name string, attachmentPath string) error {

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
      <p style="color: black;">Thank you for choosing <strong>Mind Sarthi</strong>. We wish you the best in your journey of self-discovery!</p>
	  <p style="color: black;">We would love to hear your feedback on your experience. Please feel free to reply to this email and share your thoughts with us.</p>
      <p style="color: black;">Best regards,<br>Nitish</p>
    `, name)

	err := sendEmail(to, "Insights Unlocked: Your BIG 5 Personality Assessment Report is Ready!", htmlBody, attachmentPath)

	return err
}

func SendBIG5ReportWithLink(to string, name string, link string) error {

	htmlBody := fmt.Sprintf(`
      <p style="color: black;">Dear %s,</p>
      <p style="color: black;">We hope this email finds you well. Thank you for completing the Big 5 Personality Test! We have attached your personalized report, which provides detailed insights into your personality.</p>
      <p style="color: black;"><strong>What’s inside your report:</strong></p>
      <ul style="color: black;">
        <li>A breakdown of your scores for each of the five personality traits</li>
        <li>Personalized insights based on your responses</li>
        <li>Practical tips for personal growth and development</li>
      </ul>
      <p style="color: black;">You can access your report <a href="%s">here</a>.</p>
      <p style="color: black;">By understanding your personality, you can gain valuable insights that can enhance your personal and professional life.</p>
      <p style="color: black;">Thank you for choosing <strong>Mind Sarthi</strong>. We wish you the best in your journey of self-discovery!</p>
	  <p style="color: black;">We would love to hear your feedback on your experience. Please feel free to reply to this email and share your thoughts with us.</p>
      <p style="color: black;">Best regards,<br>Nitish</p>
    `, name, link)

	err := sendEmail(to, "Insights Unlocked: Your BIG 5 Personality Assessment Report is Ready!", htmlBody, attachmentPath)

	return err
}

func Mail() {

	pdfFile := "report.pdf"

	recipientEmail := "blah@gmail.com"
	subject := "Your Report"
	body := "Hey, here's your report!"
	sendEmail(recipientEmail, subject, body, pdfFile)
}
