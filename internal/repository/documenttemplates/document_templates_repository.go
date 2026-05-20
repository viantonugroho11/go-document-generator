package documenttemplates

import (
	"context"

	tplEntity "go-document-generator/internal/entity/documenttemplates"
	"go-document-generator/internal/entity/enums"
	"go-document-generator/internal/shared/pagination"

	"gorm.io/gorm"
)

type ListFilter struct {
	TenantID *string
	Code     string
	Category string
	IsActive *bool
	Engine   enums.TemplateEngine
	Page     pagination.Params
}

type DocumentTemplatesRepository interface {
	Create(ctx context.Context, tx *gorm.DB, t tplEntity.Template) (tplEntity.Template, error)
	GetByID(ctx context.Context, tx *gorm.DB, id int64, tenantID *string) (tplEntity.Template, error)
	GetByCode(ctx context.Context, tx *gorm.DB, code string, tenantID *string) (tplEntity.Template, error)
	List(ctx context.Context, tx *gorm.DB, f ListFilter) ([]tplEntity.Template, int64, error)
	Update(ctx context.Context, tx *gorm.DB, t tplEntity.Template) (tplEntity.Template, error)
	Deactivate(ctx context.Context, tx *gorm.DB, id int64, tenantID *string, updatedBy *string) error
}
