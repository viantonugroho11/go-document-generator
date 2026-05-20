package documenttemplateversions

import (
	"context"

	verEntity "go-document-generator/internal/entity/documenttemplateversions"

	"gorm.io/gorm"
)

type DocumentTemplateVersionsRepository interface {
	Create(ctx context.Context, tx *gorm.DB, v verEntity.TemplateVersion) (verEntity.TemplateVersion, error)
	GetByID(ctx context.Context, tx *gorm.DB, templateID, versionID int64, tenantID *string) (verEntity.TemplateVersion, error)
	ListByTemplateID(ctx context.Context, tx *gorm.DB, templateID int64, tenantID *string, isPublished *bool) ([]verEntity.TemplateVersion, error)
	GetLatestPublished(ctx context.Context, tx *gorm.DB, templateID int64, tenantID *string) (verEntity.TemplateVersion, error)
	GetByTemplateAndVersion(ctx context.Context, tx *gorm.DB, templateID int64, version int, tenantID *string) (verEntity.TemplateVersion, error)
	NextVersionNumber(ctx context.Context, tx *gorm.DB, templateID int64) (int, error)
	Publish(ctx context.Context, tx *gorm.DB, templateID, versionID int64, tenantID *string) (verEntity.TemplateVersion, error)
	UnpublishOthers(ctx context.Context, tx *gorm.DB, templateID, exceptVersionID int64) error
}
