package documentrenderlogs

import (
	"context"

	logEntity "go-document-generator/internal/entity/documentrenderlogs"
	"go-document-generator/internal/shared/pagination"

	"gorm.io/gorm"
)

type DocumentRenderLogsRepository interface {
	Create(ctx context.Context, tx *gorm.DB, l logEntity.RenderLog) (logEntity.RenderLog, error)
	ListByDocumentID(ctx context.Context, tx *gorm.DB, documentID int64, page pagination.Params) ([]logEntity.RenderLog, int64, error)
}
