package postgres

import (
	"context"
	"errors"
	"time"

	entity "go-document-generator/internal/entity/documents"
	repo "go-document-generator/internal/repository/documents"
	"go-document-generator/internal/repository/documents/model"

	"gorm.io/gorm"
)

type documentsRepository struct {
	db *gorm.DB
}

func NewDocumentsRepository(db *gorm.DB) repo.DocumentsRepository {
	return &documentsRepository{db: db}
}

func (r *documentsRepository) Create(ctx context.Context, doc entity.Document) (entity.Document, error) {
	m := model.FromEntity(doc)
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = m.CreatedAt
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return entity.Document{}, err
	}
	return m.ToEntity(), nil
}

func (r *documentsRepository) GetByID(ctx context.Context, id int64) (entity.Document, error) {
	var m model.Document
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Document{}, errors.New("document not found")
	}
	return m.ToEntity(), err
}

func (r *documentsRepository) List(ctx context.Context) ([]entity.Document, error) {
	var rows []model.Document
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]entity.Document, 0, len(rows))
	for _, m := range rows {
		result = append(result, m.ToEntity())
	}
	return result, nil
}

func (r *documentsRepository) Update(ctx context.Context, doc entity.Document) (entity.Document, error) {
	updates := model.FromEntity(doc)
	tx := r.db.WithContext(ctx).Model(&model.Document{}).Where("id = ?", doc.ID).Updates(updates)
	if tx.Error != nil {
		return entity.Document{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return entity.Document{}, errors.New("document not found")
	}
	return updates.ToEntity(), nil
}

func (r *documentsRepository) Delete(ctx context.Context, id int64) error {
	tx := r.db.WithContext(ctx).Delete(&model.Document{}, "id = ?", id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("document not found")
	}
	return nil
}


