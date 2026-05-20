package documents

import (
	"context"
	"time"

	docEntity "go-document-generator/internal/entity/documents"
	"go-document-generator/internal/entity/enums"
	"go-document-generator/internal/shared/pagination"

	"gorm.io/gorm"
)

type ListFilter struct {
	TenantID     *string
	RequestID    string
	Status       enums.DocumentStatus
	TemplateCode string
	DmsStatus    enums.DmsStatus
	CallbackStatus enums.CallbackStatus
	CreatedFrom  *time.Time
	CreatedTo    *time.Time
	Page         pagination.Params
}

type DocumentsRepository interface {
	Create(ctx context.Context, tx *gorm.DB, d docEntity.Document) (docEntity.Document, error)
	GetByID(ctx context.Context, tx *gorm.DB, id int64, tenantID *string) (docEntity.Document, error)
	GetByRequestID(ctx context.Context, tx *gorm.DB, requestID string, tenantID *string) (docEntity.Document, error)
	List(ctx context.Context, tx *gorm.DB, f ListFilter) ([]docEntity.Document, int64, error)
	Update(ctx context.Context, tx *gorm.DB, d docEntity.Document) (docEntity.Document, error)
	SoftDelete(ctx context.Context, tx *gorm.DB, id int64, tenantID *string) error
}
