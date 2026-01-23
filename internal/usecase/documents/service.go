package documents

import (
	"context"
	"fmt"

	entityDocuments "go-document-generator/internal/entity/documents"
	entityTemplates "go-document-generator/internal/entity/documenttemplates"
	entityDocumentVersions "go-document-generator/internal/entity/documenttemplateversions"
	"go-document-generator/internal/shared/validators"
)

type (
	Generator interface {
		Generate(ctx context.Context, templateSource string, data any) ([]byte, string, error)
	}

	GeneratorSelector interface {
		Select(outputFormat string, engine string) Generator
	}

	Usecase interface {
		GenerateByVersionID(
			ctx context.Context,
			versionID int64,
			payload map[string]any,
		) (bytes []byte, contentType string, tmpl entityTemplates.DocumentTemplate, err error)
	}

	docsGetter interface {
		GetByID(ctx context.Context, id int64) (entityDocuments.Document, error)
	}

	docsSaver interface {
		Save(ctx context.Context, doc entityDocuments.Document) (entityDocuments.Document, error)
	}

	tmplRepo interface {
		GetByID(ctx context.Context, id int64) (entityTemplates.DocumentTemplate, error)
		GetByCode(ctx context.Context, code string) (entityTemplates.DocumentTemplate, error)
	}

	verRepo interface {
		GetByID(ctx context.Context, id int64) (entityDocumentVersions.DocumentTemplateVersion, error)
	}

	publisher interface {
		Publish(ctx context.Context, doc entityDocuments.Document) error
	}
)

// Service mengorkestrasi pengambilan template/version dan pembangkitan dokumen.
type service struct {
	docsGetter  docsGetter
	docsSaver   docsSaver
	tmplRepo    tmplRepo
	verRepo     verRepo
	genSelector GeneratorSelector
	publisher   publisher
}

// NewService membuat implementasi konkret Usecase (mengembalikan *Service untuk fleksibilitas).
func NewService(
	docsGetter docsGetter,
	docsSaver docsSaver,
	tmplRepo tmplRepo,
	verRepo verRepo,
	genSelector GeneratorSelector,
	publisher publisher,
) Usecase {
	return &service{
		docsGetter:  docsGetter,
		docsSaver:   docsSaver,
		tmplRepo:    tmplRepo,
		verRepo:     verRepo,
		genSelector: genSelector,
		publisher:   publisher,
	}
}

// GenerateByVersionID menghasilkan dokumen berdasarkan ID versi template.
// Mengambil template untuk mengetahui output_format dan engine.
func (s *service) GenerateByVersionID(
	ctx context.Context,
	versionID int64,
	payload map[string]any,
) (bytes []byte, contentType string, tmpl entityTemplates.DocumentTemplate, err error) {
	ver, err := s.verRepo.GetByID(ctx, versionID)
	if err != nil {
		return nil, "", entityTemplates.DocumentTemplate{}, err
	}
	tmpl, err = s.tmplRepo.GetByID(ctx, ver.TemplateID)
	if err != nil {
		return nil, "", entityTemplates.DocumentTemplate{}, err
	}
	gen := s.genSelector.Select(tmpl.OutputFormat, tmpl.Engine)
	if gen == nil {
		return nil, "", entityTemplates.DocumentTemplate{}, fmt.Errorf("no generator for format=%s engine=%s", tmpl.OutputFormat, tmpl.Engine)
	}
	out, ct, err := gen.Generate(ctx, ver.Content, payload)
	if err != nil {
		return nil, "", entityTemplates.DocumentTemplate{}, err
	}
	return out, ct, tmpl, nil
}

// CreateDocument creates a new document.
func (s *service) Create(ctx context.Context, doc entityDocuments.Document) (entityDocuments.Document, error) {

	ver, err := s.verRepo.GetByID(ctx, int64(*doc.TemplateVersion))
	if err != nil {
		return entityDocuments.Document{}, err
	}
	_, err = s.tmplRepo.GetByCode(ctx, doc.TemplateCode)
	if err != nil {
		return entityDocuments.Document{}, err
	}

	// validate schema template
	if err = validators.ValidateSchema(ver.Schema, doc.Payload); err != nil {
		return entityDocuments.Document{}, err
	}

	// save document
	doc, err = s.docsSaver.Save(ctx, doc)
	if err != nil {
		return entityDocuments.Document{}, err
	}

	// publish document
	err = s.publisher.Publish(ctx, doc)
	if err != nil {
		return entityDocuments.Document{}, err
	}

	return doc, nil
}
