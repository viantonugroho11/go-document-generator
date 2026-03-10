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
	m := model.FromEntity(v)
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return entity.DocumentTemplateVersion{}, err
	}
	return m.ToEntity(), nil
}

func (r *documentTemplateVersionsRepository) GetByID(ctx context.Context, id int64) (entity.DocumentTemplateVersion, error) {
	var m model.DocumentTemplateVersion
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.DocumentTemplateVersion{}, errors.New("document template version not found")
	}
	return m.ToEntity(), err
}

func (r *documentTemplateVersionsRepository) List(ctx context.Context) ([]entity.DocumentTemplateVersion, error) {
	var rows []model.DocumentTemplateVersion
	if err := r.db.WithContext(ctx).Order("template_id ASC, version DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	res := make([]entity.DocumentTemplateVersion, 0, len(rows))
	for _, m := range rows {
		res = append(res, m.ToEntity())
	}
	return res, nil
}

func (r *documentTemplateVersionsRepository) ListByTemplateID(ctx context.Context, templateID int64) ([]entity.DocumentTemplateVersion, error) {
	var rows []model.DocumentTemplateVersion
	if err := r.db.WithContext(ctx).
		Where("template_id = ?", templateID).
		Order("version DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]entity.DocumentTemplateVersion, 0, len(rows))
	for _, m := range rows {
		result = append(result, m.ToEntity())
	}
	return result, nil
}

func (r *documentTemplateVersionsRepository) GetLatestVersionNumber(ctx context.Context, templateID int64) (int, error) {
	var maxVersion int
	err := r.db.WithContext(ctx).
		Model(&model.DocumentTemplateVersion{}).
		Select("COALESCE(MAX(version), 0)").
		Where("template_id = ?", templateID).
		Scan(&maxVersion).Error
	if err != nil {
		return 0, err
	}
	return maxVersion, nil
}

func (r *documentTemplateVersionsRepository) Update(ctx context.Context, v entity.DocumentTemplateVersion) (entity.DocumentTemplateVersion, error) {
	updates := model.FromEntity(v)
	tx := r.db.WithContext(ctx).Model(&model.DocumentTemplateVersion{}).Where("id = ?", v.ID).Updates(updates)
	if tx.Error != nil {
		return entity.DocumentTemplateVersion{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return entity.DocumentTemplateVersion{}, errors.New("document template version not found")
	}
	return updates.ToEntity(), nil
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
