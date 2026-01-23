package documents

import (
	"strings"

	"go-document-generator/internal/infrastructure/documents/csv"
	"go-document-generator/internal/infrastructure/documents/pdf"
)

// NewGeneratorFromTemplate mengembalikan Generator berdasarkan outputFormat dan engine template.
// - outputFormat: "pdf" atau "csv"
// - engineType: templating.EngineType (tmpl/html/...)
func NewGeneratorFromTemplate(outputFormat string, engineType any) Generator {
	switch strings.ToLower(outputFormat) {
	case "pdf":
		return pdf.NewWKHTMLToPDFGenerator()
	case "csv":
		return csv.NewCSVGenerator()
	default:
		// fallback: CSV generator
		return csv.NewCSVGenerator()
	}
}
