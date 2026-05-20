package postgres

import (
	"context"
	"errors"
	"time"

	verEntity "go-document-generator/internal/entity/documenttemplateversions"
	repo "go-document-generator/internal/repository/documenttemplateversions"
	"go-document-generator/internal/repository/documenttemplateversions/model"
	"go-document-generator/internal/shared/apperror"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewDocumentTemplateVersionsRepository(db *gorm.DB) repo.DocumentTemplateVersionsRepository {
	return &repository{db: db}
}

func (r *repository) conn(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *repository) Create(ctx context.Context, tx *gorm.DB, v verEntity.TemplateVersion) (verEntity.TemplateVersion, error) {
	m := model.ToModel(v)
	m.CreatedAt = time.Now().UTC()
	if err := r.conn(tx).WithContext(ctx).Create(&m).Error; err != nil {
		return verEntity.TemplateVersion{}, err
	}
	return model.ToEntity(&m), nil
}

func (r *repository) GetByID(ctx context.Context, tx *gorm.DB, templateID, versionID int64, tenantID *string) (verEntity.TemplateVersion, error) {
	var m model.DocumentTemplateVersion
	q := r.conn(tx).WithContext(ctx).Where("id = ? AND template_id = ?", versionID, templateID)
	if tenantID != nil {
		q = q.Where("tenant_id = ?", *tenantID)
	}
	if err := q.First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return verEntity.TemplateVersion{}, apperror.ErrNotFound
		}
		return verEntity.TemplateVersion{}, err
	}
	return model.ToEntity(&m), nil
}

func (r *repository) ListByTemplateID(ctx context.Context, tx *gorm.DB, templateID int64, tenantID *string, isPublished *bool) ([]verEntity.TemplateVersion, error) {
	q := r.conn(tx).WithContext(ctx).Where("template_id = ?", templateID)
	if tenantID != nil {
		q = q.Where("tenant_id = ?", *tenantID)
	}
	if isPublished != nil {
		q = q.Where("is_published = ?", *isPublished)
	}
	q = q.Order("version DESC")

	var rows []model.DocumentTemplateVersion
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]verEntity.TemplateVersion, len(rows))
	for i := range rows {
		out[i] = model.ToEntity(&rows[i])
	}
	return out, nil
}

func (r *repository) GetLatestPublished(ctx context.Context, tx *gorm.DB, templateID int64, tenantID *string) (verEntity.TemplateVersion, error) {
	published := true
	rows, err := r.ListByTemplateID(ctx, tx, templateID, tenantID, &published)
	if err != nil {
		return verEntity.TemplateVersion{}, err
	}
	if len(rows) == 0 {
		return verEntity.TemplateVersion{}, apperror.ErrNotFound
	}
	return rows[0], nil
}

func (r *repository) GetByTemplateAndVersion(ctx context.Context, tx *gorm.DB, templateID int64, version int, tenantID *string) (verEntity.TemplateVersion, error) {
	var m model.DocumentTemplateVersion
	q := r.conn(tx).WithContext(ctx).Where("template_id = ? AND version = ?", templateID, version)
	if tenantID != nil {
		q = q.Where("tenant_id = ?", *tenantID)
	}
	if err := q.First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return verEntity.TemplateVersion{}, apperror.ErrNotFound
		}
		return verEntity.TemplateVersion{}, err
	}
	return model.ToEntity(&m), nil
}

func (r *repository) NextVersionNumber(ctx context.Context, tx *gorm.DB, templateID int64) (int, error) {
	var maxVersion int
	err := r.conn(tx).WithContext(ctx).
		Model(&model.DocumentTemplateVersion{}).
		Where("template_id = ?", templateID).
		Select("COALESCE(MAX(version), 0)").
		Scan(&maxVersion).Error
	if err != nil {
		return 0, err
	}
	return maxVersion + 1, nil
}

func (r *repository) UnpublishOthers(ctx context.Context, tx *gorm.DB, templateID, exceptVersionID int64) error {
	return r.conn(tx).WithContext(ctx).
		Model(&model.DocumentTemplateVersion{}).
		Where("template_id = ? AND id <> ?", templateID, exceptVersionID).
		Where("is_published = ?", true).
		Updates(map[string]any{
			"is_published": false,
			"published_at": nil,
		}).Error
}

func (r *repository) Publish(ctx context.Context, tx *gorm.DB, templateID, versionID int64, tenantID *string) (verEntity.TemplateVersion, error) {
	now := time.Now().UTC()
	q := r.conn(tx).WithContext(ctx).
		Model(&model.DocumentTemplateVersion{}).
		Where("id = ? AND template_id = ?", versionID, templateID)
	if tenantID != nil {
		q = q.Where("tenant_id = ?", *tenantID)
	}
	res := q.Updates(map[string]any{
		"is_published": true,
		"published_at": now,
	})
	if res.Error != nil {
		return verEntity.TemplateVersion{}, res.Error
	}
	if res.RowsAffected == 0 {
		return verEntity.TemplateVersion{}, apperror.ErrNotFound
	}
	return r.GetByID(ctx, tx, templateID, versionID, tenantID)
}
