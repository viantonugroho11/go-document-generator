package documents

import (
	"strings"

	"go-document-generator/internal/infrastructure/documents/csv"
	"go-document-generator/internal/infrastructure/documents/html"
	"go-document-generator/internal/infrastructure/documents/pdf"
	usecasedoc "go-document-generator/internal/usecase/documents"
)

// Selector mengimplementasikan usecase GeneratorSelector.
type Selector struct{}

func NewSelector() *Selector { return &Selector{} }

func (s *Selector) Select(outputFormat string, engine string) usecasedoc.Generator {
	switch strings.ToUpper(outputFormat) {
	case "PDF":
		return pdf.NewWKHTMLToPDFGenerator()
	case "HTML":
		return html.NewGenerator()
	case "DOCX":
		// DOCX belum diimplementasikan; fallback ke HTML template output.
		return html.NewGenerator()
	default:
		switch strings.ToUpper(engine) {
		case "HTML":
			return html.NewGenerator()
		default:
			return csv.NewCSVGenerator()
		}
	}
}
