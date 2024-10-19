package API

import (
	"github.com/jung-kurt/gofpdf/v2"
)

func main() {
	// Create a new PDF instance
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set font for the header (Extraversion)
	pdf.SetFont("Arial", "B", 40)
	pdf.SetTextColor(50, 50, 50)
	pdf.Cell(0, 20, "EXTRAVERSION")
	pdf.Ln(20)

	// Insert the logo in the top right corner
	pdf.ImageOptions("2.png", 160, 10, 40, 40, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	// Set footer font and text
	pdf.SetY(-15)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(0, 153, 204)
	footerText := "CONFIDENTIALITY NOTICE: This report contains sensitive personal information and is intended solely for the use of the individual named above. Unauthorized disclosure or reproduction of this document is strictly prohibited."
	pdf.MultiCell(0, 10, footerText, "", "C", false)

	// Save the PDF to a file
	err := pdf.OutputFileAndClose("extraversion_report.pdf")
	if err != nil {
		panic(err)
	}
}
