package API

import (
	"github.com/jung-kurt/gofpdf"
	"log"
)

// Struct to hold the AI output for each domain
type DomainContent struct {
	Introduction     string
	CareerAcademia   string
	Relationship     string
	StrengthWeakness string
}

// Function to generate the PDF with 5 pages for Big Five domains
// Modify to accept filename as a parameter
func GenerateBigFivePDF(contents map[string]DomainContent, filename string) error {

	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.SetFont("Arial", "", 14)

	// Define domain order
	domains := []string{"Neuroticism", "Extraversion", "Openness", "Agreeableness", "Conscientiousness"}

	// Loop through the domains and add each one to a new page
	for _, domain := range domains {
		// Add a new page for each domain
		pdf.AddPage()

		// Get the content for the current domain
		content := contents[domain]

		// Add title (domain name)
		pdf.SetFont("Arial", "B", 16)
		pdf.Cell(40, 10, domain)
		pdf.Ln(12)

		// Add the sections for Introduction, Career & Academia, etc.
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(40, 10, "Introduction:")
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 12)
		pdf.MultiCell(190, 10, content.Introduction, "", "", false)
		pdf.Ln(5)

		// Career & Academia
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(40, 10, "Career & Academia:")
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 12)
		pdf.MultiCell(190, 10, content.CareerAcademia, "", "", false)
		pdf.Ln(5)

		// Relationship
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(40, 10, "Relationship:")
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 12)
		pdf.MultiCell(190, 10, content.Relationship, "", "", false)
		pdf.Ln(5)

		// Strength & Weakness
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(40, 10, "Strength & Weakness:")
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 12)
		pdf.MultiCell(190, 10, content.StrengthWeakness, "", "", false)
		pdf.Ln(5)

	}

	// Use TestId as the PDF filename
	pdfFilename := filename + ".pdf"
	err := pdf.OutputFileAndClose(pdfFilename)
	if err != nil {
		return err
	}
	return nil
}

func CreatePDF() {
	// Example content for each domain (this will be replaced by AI output in real use)
	contents := map[string]DomainContent{
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
