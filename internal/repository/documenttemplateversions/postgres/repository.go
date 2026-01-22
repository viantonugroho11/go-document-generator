package postgres

import (
	"context"
	"errors"
	"time"

	entity "go-document-generator/internal/entity/documenttemplateversions"
	repo "go-document-generator/internal/repository/documenttemplateversions"
	"go-document-generator/internal/repository/documenttemplateversions/model"
	"gorm.io/gorm"
)

type documentTemplateVersionsRepository struct {
	db *gorm.DB
}

func NewDocumentTemplateVersionsRepository(db *gorm.DB) repo.DocumentTemplateVersionsRepository {
	return &documentTemplateVersionsRepository{db: db}
}

func (r *documentTemplateVersionsRepository) Create(ctx context.Context, v entity.DocumentTemplateVersion) (entity.DocumentTemplateVersion, error) {
	m := toModel(v)
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return entity.DocumentTemplateVersion{}, err
	}
	return toEntity(m), nil
}

func (r *documentTemplateVersionsRepository) GetByID(ctx context.Context, id int64) (entity.DocumentTemplateVersion, error) {
	var m model.DocumentTemplateVersion
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.DocumentTemplateVersion{}, errors.New("document template version not found")
	}
	return toEntity(m), err
}

func (r *documentTemplateVersionsRepository) List(ctx context.Context) ([]entity.DocumentTemplateVersion, error) {
	var rows []model.DocumentTemplateVersion
	if err := r.db.WithContext(ctx).Order("template_id ASC, version DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	res := make([]entity.DocumentTemplateVersion, 0, len(rows))
	for _, m := range rows {
		res = append(res, toEntity(m))
	}
	return res, nil
}

func (r *documentTemplateVersionsRepository) Update(ctx context.Context, v entity.DocumentTemplateVersion) (entity.DocumentTemplateVersion, error) {
	updates := map[string]any{
		"template_id":    v.TemplateID,
		"version":        v.Version,
		"content":        v.Content,
		"schema":         v.Schema,
		"sample_payload": v.SamplePayload,
		"is_published":   v.IsPublished,
		"published_at":   v.PublishedAt,
	}
	tx := r.db.WithContext(ctx).Model(&model.DocumentTemplateVersion{}).Where("id = ?", v.ID).Updates(updates)
	if tx.Error != nil {
		return entity.DocumentTemplateVersion{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return entity.DocumentTemplateVersion{}, errors.New("document template version not found")
	}
	return v, nil
}

func (r *documentTemplateVersionsRepository) Delete(ctx context.Context, id int64) error {
	tx := r.db.WithContext(ctx).Delete(&model.DocumentTemplateVersion{}, "id = ?", id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("document template version not found")
	}
	return nil
}

func toModel(e entity.DocumentTemplateVersion) model.DocumentTemplateVersion {
	return model.DocumentTemplateVersion{
		ID:            e.ID,
		TemplateID:    e.TemplateID,
		Version:       e.Version,
		Content:       e.Content,
		Schema:        e.Schema,
		SamplePayload: e.SamplePayload,
		IsPublished:   e.IsPublished,
		PublishedAt:   e.PublishedAt,
		CreatedAt:     e.CreatedAt,
	}
}

func toEntity(m model.DocumentTemplateVersion) entity.DocumentTemplateVersion {
	return entity.DocumentTemplateVersion{
		ID:            m.ID,
		TemplateID:    m.TemplateID,
		Version:       m.Version,
		Content:       m.Content,
		Schema:        m.Schema,
		SamplePayload: m.SamplePayload,
		IsPublished:   m.IsPublished,
		PublishedAt:   m.PublishedAt,
		CreatedAt:     m.CreatedAt,
	}
}

