package lib

import (
	"github.com/go-pdf/fpdf"

	"log"
)

// CreatePdf creates a simple PDF with "Hello, world"
func CreatePdf() {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, world")
	err := pdf.OutputFileAndClose("hello.pdf")
	if err != nil {
		panic(err) // handle error if file generation fails
	}
}

func CreatePdfWithBg() {
	// Create a new PDF in portrait mode, A4 size
	pdf := fpdf.New("P", "mm", "A4", "")

	// Set PDF dimensions (A4 size)
	width, height := 210.0, 297.0

	// Add a page to the PDF
	pdf.AddPage()

	// Provide the correct image path (with double backslashes or forward slashes)
	imgPath := "C:\\Users\\Rishi Raj Ganguly\\cognify-api-gateway\\pdfBg.png" // Double backslashes for Windows path
	// Alternatively, you can use forward slashes:
	// imgPath := "C:/Users/Rishi Raj Ganguly/cognify-api-gateway/pdfBg.png"

	// Use ImageOptions to define the background behavior (scaling and positioning)
	opt := fpdf.ImageOptions{
		ImageType: "PNG", // Define the image type (PNG in this case)
		ReadDpi:   true,  // Use the DPI from the image
	}

	// Insert the background image and scale it to fit the entire page
	pdf.ImageOptions(imgPath, 0, 0, width, height, false, opt, 0, "")

	// // Add some text on top of the image background
	// pdf.SetFont("Arial", "B", 16)
	// pdf.SetTextColor(255, 255, 255) // White text
	// pdf.Cell(40, 10, "Hello, world on a background image!")

	// Output the PDF to a file
	err := pdf.OutputFileAndClose("output_with_background.pdf")
	if err != nil {
		log.Fatalf("Failed to generate PDF: %s", err)
	}
}
