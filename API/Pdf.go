package API

import (
	"github.com/jung-kurt/gofpdf"
	"log"
)

func generateBig5Report(data []Domain, narrative map[string]string, outputPath string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle("Big 5 Personality Assessment Report", false)

	for _, domain := range data {
		pdf.AddPage()
		pdf.SetFont("Arial", "B", 16)
		pdf.Cell(0, 10, domain.Name+" Assessment")
		pdf.Ln(12)

		// Print subdomains
		pdf.SetFont("Arial", "", 12)
		for _, sub := range domain.Subdomains {
			pdf.Cell(0, 10, sub.Name+": "+sub.Score)
			pdf.Ln(8)
		}

		// Add the narrative content
		pdf.Ln(12)
		pdf.MultiCell(0, 10, narrative[domain.Name], "", "", false)
	}

	err := pdf.OutputFileAndClose(outputPath)
	if err != nil {
		return err
	}

	log.Printf("PDF successfully generated at: %s", outputPath)
	return nil
}
