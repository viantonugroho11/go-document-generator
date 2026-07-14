package transitions

import (
	"context"
	"time"

	"go-document-generator/internal/entity/enums"
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

// StorageProvider abstraksi penyimpanan file (sama signature dengan shared/storage.Provider).
type StorageProvider interface {
	Save(ctx context.Context, documentID int64, requestID, ext string, data []byte) (path, fileName string, err error)
	PresignedURL(ctx context.Context, path string, ttl time.Duration) (string, error)
	ProviderName() enums.StorageProvider
}

// Deps dependensi untuk handler transisi status dokumen.
type Deps struct {
	Templates tplrepo.DocumentTemplatesRepository
	Versions  verrepo.DocumentTemplateVersionsRepository
	Selector  GeneratorSelector
	Storage   StorageProvider
}
