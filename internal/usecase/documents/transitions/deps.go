package transitions

import (
	"context"

	tplrepo "go-document-generator/internal/repository/documenttemplates"
	verrepo "go-document-generator/internal/repository/documenttemplateversions"
)

// Generator merender dokumen (kontrak sama dengan documents.Generator).
type Generator interface {
	Generate(ctx context.Context, templateSource string, data any) ([]byte, string, error)
}

// GeneratorSelector memilih engine render.
type GeneratorSelector interface {
	Select(outputFormat string, engine string) Generator
}

// Deps dependensi untuk handler transisi status dokumen.
type Deps struct {
	Templates tplrepo.DocumentTemplatesRepository
	Versions  verrepo.DocumentTemplateVersionsRepository
	Selector  GeneratorSelector
}
