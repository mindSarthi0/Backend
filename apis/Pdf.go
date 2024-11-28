package API

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"myproject/constants"
)

// Struct to hold the AI output for each domain

// Function to generate the PDF with 5 pages for Big Five domains
// Modify to accept filename as a parameter
func GenerateBigFivePDF(contents map[string]string, name string, filename string) error {

	// Create a new PDF instance
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10) // Set margins

	// Add the cover page
	pdf.AddPage()

	// fmt.Println(":::::::::Contents:::::::::", contents)

	// Path to the cover image
	coverImgPath := "./Reports/Template/co.png"

	// Define image options for cover
	coverOpt := gofpdf.ImageOptions{
		ImageType: "PNG",
		ReadDpi:   true,
	}

	// Insert cover image to fit the entire page
	width, height := 210.0, 297.0 // A4 dimensions
	pdf.ImageOptions(coverImgPath, 0, 0, width, height, false, coverOpt, 0, "")
	pdf.SetY(108.5)
	pdf.SetX(60)
	pdf.SetFont("Arial", "B", 24)
	pdf.SetTextColor(17, 45, 78)
	if name != "" {
		pdf.Cell(40, 10, name) // Heading in small caps
		pdf.Ln(8)
	}

	// Define the domain order for the next pages
	// Loop through the domains and add each one to a new page
	for key, value := range constants.BIG_5_Report {
		// Add a new page for each domain
		pdf.AddPage()
		if err := pdf.Error(); err != nil {
			fmt.Println("failed to add a new page: %w", err)
		}

		// Construct the correct image path (use forward slashes for Go)
		imgPath := fmt.Sprintf("./Reports/Template/%s.png", value)

		// Define image options (e.g., scaling and positioning)
		opt := gofpdf.ImageOptions{
			ImageType: "PNG",
			ReadDpi:   true,
		}

		// Insert the image; for the whole page, set width to 210mm (A4 width) and height to 297mm (A4 height)
		pdf.ImageOptions(imgPath, 0, 0, width, height, false, opt, 0, "")

		if err := pdf.Error(); err != nil {
			fmt.Println("failed to add image: %w", err)
		}

		// Get the content for the current domain
		content, ok := contents[key]

		if !ok {
			continue // If the domain doesn't exist in the content map, skip it
		}

		// Set some Y offset after the image to start the text content
		pdf.SetY(50) // Adjust based on your needs
		if err := pdf.Error(); err != nil {
			fmt.Println("failed to set Y offset: %w", err)
		}

		// log.Println(content)

		addContentSection(pdf, "", content, 24)
		if err := pdf.Error(); err != nil {
			fmt.Println("failed to add page: %w", err)
		}
	}

	// Save the PDF to a file
	pdfFilename := filename + ".pdf"
	err := pdf.OutputFileAndClose(pdfFilename)

	if err != nil {
		println("PDF generation error", err)
		return err
	}
	return nil
}

// Helper function to add content sections to the PDF
// Adds the title in small caps, with specific heading and body fonts
func addContentSection(pdf *gofpdf.Fpdf, title, content string, fontSize int) {
	// Set the heading font to Calibri bold, small caps
	pdf.SetFont("Arial", "B", float64(fontSize))
	pdf.SetTextColor(17, 45, 78)
	if title != "" {
		pdf.Cell(40, 10, title+":") // Heading in small caps
		pdf.Ln(8)
	}

	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(17, 45, 78)
	pdf.MultiCell(190, 10, content, "", "", false)
	pdf.Ln(10) // Additional gap after the body text
}

// Function to create the PDF
func CreatePDF() {
	// Example content for each domain (this will be replaced by AI output in real use)

	// // Provide a filename for the generated PDF
	// filename := "BigFiveReport"

	// // Generate the PDF with the filename passed as a second argument
	// err := GenerateBigFivePDF(contents, "User name", filename)
	// if err != nil {
	// 	log.Fatalf("Failed to generate PDF: %v", err)
	// }
}
