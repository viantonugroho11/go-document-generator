package documents

import (
	"context"

	entity "go-document-generator/internal/entity/documents"
)

// Interface repository untuk entity Document.
// Implementasi (Postgres/Mongo/dll) harus memenuhi kontrak ini.
type DocumentsRepository interface {
	Create(ctx context.Context, doc entity.Document) (entity.Document, error)
	GetByID(ctx context.Context, id int64) (entity.Document, error)
	List(ctx context.Context) ([]entity.Document, error)
	Update(ctx context.Context, doc entity.Document) (entity.Document, error)
	Delete(ctx context.Context, id int64) error
}

