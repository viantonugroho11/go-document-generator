package documents

import (
	"strings"

	"go-document-generator/internal/infrastructure/documents/csv"
	"go-document-generator/internal/infrastructure/documents/pdf"
	usecasedoc "go-document-generator/internal/usecase/documents"
)

// Selector mengimplementasikan usecase GeneratorSelector.
type Selector struct{}

func NewSelector() *Selector { return &Selector{} }

func (s *Selector) Select(outputFormat string, engine string) usecasedoc.Generator {
	switch strings.ToLower(outputFormat) {
	case "pdf":
		return pdf.NewWKHTMLToPDFGenerator()
	case "csv":
		return csv.NewCSVGenerator()
	default:
		return csv.NewCSVGenerator()
	}
}
