package documenttemplates

import (
	"context"

	entity "go-document-generator/internal/entity/documenttemplates"
)

// Interface repository untuk entity DocumentTemplate.
type DocumentTemplatesRepository interface {
	Create(ctx context.Context, tmpl entity.DocumentTemplate) (entity.DocumentTemplate, error)
	GetByID(ctx context.Context, id int64) (entity.DocumentTemplate, error)
	List(ctx context.Context) ([]entity.DocumentTemplate, error)
	Update(ctx context.Context, tmpl entity.DocumentTemplate) (entity.DocumentTemplate, error)
	Delete(ctx context.Context, id int64) error
}

