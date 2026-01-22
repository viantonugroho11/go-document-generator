package documenttemplateversions

import (
	"context"

	entity "go-document-generator/internal/entity/documenttemplateversions"
)

// Interface repository untuk entity DocumentTemplateVersion.
type DocumentTemplateVersionsRepository interface {
	Create(ctx context.Context, v entity.DocumentTemplateVersion) (entity.DocumentTemplateVersion, error)
	GetByID(ctx context.Context, id int64) (entity.DocumentTemplateVersion, error)
	List(ctx context.Context) ([]entity.DocumentTemplateVersion, error)
	Update(ctx context.Context, v entity.DocumentTemplateVersion) (entity.DocumentTemplateVersion, error)
	Delete(ctx context.Context, id int64) error
}

