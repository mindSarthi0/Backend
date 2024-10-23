package API

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"log"
)

// Struct to hold the AI output for each domain

// Function to generate the PDF with 5 pages for Big Five domains
// Modify to accept filename as a parameter
func GenerateBigFivePDF(contents map[string]JSONOutputFormat, filename string) error {

	// Create a new PDF instance
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10) // Set margins

	// Add the cover page
	pdf.AddPage()

	fmt.Println(":::::::::Contents:::::::::", contents)

	// Path to the cover image
	coverImgPath := "E:/Projects/Psychological Assessment/Code/cognify-api-gateway/Reports/Template/cover.png"

	// Define image options for cover
	coverOpt := gofpdf.ImageOptions{
		ImageType: "PNG",
		ReadDpi:   true,
	}

	// Insert cover image to fit the entire page
	width, height := 210.0, 297.0 // A4 dimensions
	pdf.ImageOptions(coverImgPath, 0, 0, width, height, false, coverOpt, 0, "")

	// Optional: Add text on the cover page (e.g., Report Title)
	pdf.SetY(150) // Adjust as needed
	pdf.SetFont("Arial", "B", 24)
	pdf.Cell(0, 10, "")
	pdf.Ln(12)

	pdf.SetFont("Arial", "I", 16)
	pdf.Cell(0, 10, "")
	pdf.Ln(10)

	// Define the domain order for the next pages
	domains := []string{"extraversion", "neuroticism", "openness", "agreeableness", "conscientiousness"}

	// Loop through the domains and add each one to a new page
	for _, domain := range domains {
		// Add a new page for each domain
		pdf.AddPage()

		// Construct the correct image path (use forward slashes for Go)
		imgPath := fmt.Sprintf("E:/Projects/Psychological Assessment/Code/cognify-api-gateway/Reports/Template/%s.png", domain)

		// Define image options (e.g., scaling and positioning)
		opt := gofpdf.ImageOptions{
			ImageType: "PNG",
			ReadDpi:   true,
		}

		// Insert the image; for the whole page, set width to 210mm (A4 width) and height to 297mm (A4 height)
		pdf.ImageOptions(imgPath, 0, 0, width, height, false, opt, 0, "")

		// Get the content for the current domain
		content, ok := contents[domain]
		if !ok {
			continue // If the domain doesn't exist in the content map, skip it
		}

		// Set some Y offset after the image to start the text content
		pdf.SetY(50) // Adjust based on your needs

		log.Println(content)

		// Add the sections for Introduction, Career & Academia, etc.
		addContentSection(pdf, "Introduction", content.Introduction, 24)
		addContentSection(pdf, "Career & Academia", content.CareerAcademia, 20)
		addContentSection(pdf, "Relationship", content.Relationship, 20)
		addContentSection(pdf, "Strength & Weakness", content.StrengthWeakness, 20)
	}

	// Save the PDF to a file
	pdfFilename := filename + ".pdf"
	err := pdf.OutputFileAndClose(pdfFilename)
	if err != nil {
		return err
	}
	return nil
}

// Helper function to add content sections to the PDF
// Adds the title in small caps, with specific heading and body fonts
func addContentSection(pdf *gofpdf.Fpdf, title, content string, fontSize int) {
	// Set the heading font to Calibri bold, small caps
	pdf.SetFont("Calibri", "B", float64(fontSize))
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(40, 10, title+":") // Heading in small caps
	pdf.Ln(8)

	// Set the body font to Calibri 16 pt and add space before the content
	pdf.SetFont("Calibri", "", 16)
	pdf.MultiCell(190, 10, content, "", "", false)
	pdf.Ln(10) // Additional gap after the body text
}

// Function to create the PDF
func CreatePDF() {
	// Example content for each domain (this will be replaced by AI output in real use)
	contents := map[string]JSONOutputFormat{
		"Neuroticism": {
			Introduction:     "Neuroticism reflects emotional stability and degree of negative emotions...",
			CareerAcademia:   "In a professional setting, high neuroticism may lead to stress under pressure...",
			Relationship:     "Those with high neuroticism may experience emotional turbulence in relationships...",
			StrengthWeakness: "Strength: Sensitive to emotional cues. Weakness: Prone to anxiety and mood swings.",
		},
		"Extraversion": {
			Introduction:     "Extraversion is characterized by assertiveness, social enthusiasm, and talkativeness...",
			CareerAcademia:   "Highly extroverted individuals thrive in team environments and leadership roles...",
			Relationship:     "Extroverts often have a large social circle and seek social interactions...",
			StrengthWeakness: "Strength: Strong social presence. Weakness: May dominate conversations.",
		},
		"Openness": {
			Introduction:     "Openness involves creativity, imagination, and openness to new experiences...",
			CareerAcademia:   "High openness fosters innovation and adaptability in the workplace...",
			Relationship:     "Openness can lead to deep, meaningful connections in relationships...",
			StrengthWeakness: "Strength: Creativity and curiosity. Weakness: May be perceived as unpredictable.",
		},
		"Agreeableness": {
			Introduction:     "Agreeableness reflects altruism, trust, and prosocial behavior...",
			CareerAcademia:   "Highly agreeable individuals work well in collaborative environments...",
			Relationship:     "Agreeable individuals tend to be empathetic and maintain harmonious relationships...",
			StrengthWeakness: "Strength: Strong empathy and cooperation. Weakness: May avoid conflict at personal cost.",
		},
		"Conscientiousness": {
			Introduction:     "Conscientiousness is defined by organization, dependability, and a sense of duty...",
			CareerAcademia:   "In a professional setting, conscientious individuals are reliable and goal-oriented...",
			Relationship:     "Conscientiousness leads to responsible and committed relationships...",
			StrengthWeakness: "Strength: High self-discipline. Weakness: Can be overly rigid or perfectionistic.",
		},
	}

	// Provide a filename for the generated PDF
	filename := "BigFiveReport"

	// Generate the PDF with the filename passed as a second argument
	err := GenerateBigFivePDF(contents, filename)
	if err != nil {
		log.Fatalf("Failed to generate PDF: %v", err)
	}
}
