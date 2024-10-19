package lib

import (
	"github.com/go-pdf/fpdf"
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
