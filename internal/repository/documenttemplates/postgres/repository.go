package postgres

import (
	"context"
	"errors"
	"time"

	entity "go-document-generator/internal/entity/documenttemplates"
	repo "go-document-generator/internal/repository/documenttemplates"
	"go-document-generator/internal/repository/documenttemplates/model"
	"gorm.io/gorm"
)

type documentTemplatesRepository struct {
	db *gorm.DB
}

func NewDocumentTemplatesRepository(db *gorm.DB) repo.DocumentTemplatesRepository {
	return &documentTemplatesRepository{db: db}
}

func (r *documentTemplatesRepository) Create(ctx context.Context, tmpl entity.DocumentTemplate) (entity.DocumentTemplate, error) {
	m := toModel(tmpl)
	now := time.Now()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}
	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = now
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return entity.DocumentTemplate{}, err
	}
	return toEntity(m), nil
}

func (r *documentTemplatesRepository) GetByID(ctx context.Context, id int64) (entity.DocumentTemplate, error) {
	var m model.DocumentTemplate
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.DocumentTemplate{}, errors.New("document template not found")
	}
	return toEntity(m), err
}

func (r *documentTemplatesRepository) List(ctx context.Context) ([]entity.DocumentTemplate, error) {
	var rows []model.DocumentTemplate
	if err := r.db.WithContext(ctx).Order("code ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	res := make([]entity.DocumentTemplate, 0, len(rows))
	for _, m := range rows {
		res = append(res, toEntity(m))
	}
	return res, nil
}

func (r *documentTemplatesRepository) Update(ctx context.Context, tmpl entity.DocumentTemplate) (entity.DocumentTemplate, error) {
	updates := map[string]any{
		"code":          tmpl.Code,
		"name":          tmpl.Name,
		"description":   tmpl.Description,
		"engine":        tmpl.Engine,
		"output_format": tmpl.OutputFormat,
		"is_active":     tmpl.IsActive,
		"updated_at":    time.Now(),
	}
	tx := r.db.WithContext(ctx).Model(&model.DocumentTemplate{}).Where("id = ?", tmpl.ID).Updates(updates)
	if tx.Error != nil {
		return entity.DocumentTemplate{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return entity.DocumentTemplate{}, errors.New("document template not found")
	}
	return tmpl, nil
}

func (r *documentTemplatesRepository) Delete(ctx context.Context, id int64) error {
	tx := r.db.WithContext(ctx).Delete(&model.DocumentTemplate{}, "id = ?", id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("document template not found")
	}
	return nil
}

func toModel(e entity.DocumentTemplate) model.DocumentTemplate {
	return model.DocumentTemplate{
		ID:           e.ID,
		Code:         e.Code,
		Name:         e.Name,
		Description:  e.Description,
		Engine:       e.Engine,
		OutputFormat: e.OutputFormat,
		IsActive:     e.IsActive,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

func toEntity(m model.DocumentTemplate) entity.DocumentTemplate {
	return entity.DocumentTemplate{
		ID:           m.ID,
		Code:         m.Code,
		Name:         m.Name,
		Description:  m.Description,
		Engine:       m.Engine,
		OutputFormat: m.OutputFormat,
		IsActive:     m.IsActive,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

