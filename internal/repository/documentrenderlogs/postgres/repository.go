package postgres

import (
	"context"
	"time"

	logEntity "go-document-generator/internal/entity/documentrenderlogs"
	repo "go-document-generator/internal/repository/documentrenderlogs"
	"go-document-generator/internal/repository/documentrenderlogs/model"
	"go-document-generator/internal/shared/pagination"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewDocumentRenderLogsRepository(db *gorm.DB) repo.DocumentRenderLogsRepository {
	return &repository{db: db}
}

func (r *repository) conn(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *repository) Create(ctx context.Context, tx *gorm.DB, l logEntity.RenderLog) (logEntity.RenderLog, error) {
	m := model.ToModel(l)
	m.CreatedAt = time.Now().UTC()
	if err := r.conn(tx).WithContext(ctx).Create(&m).Error; err != nil {
		return logEntity.RenderLog{}, err
	}
	return model.ToEntity(&m), nil
}

func (r *repository) ListByDocumentID(ctx context.Context, tx *gorm.DB, documentID int64, page pagination.Params) ([]logEntity.RenderLog, int64, error) {
	q := r.conn(tx).WithContext(ctx).Model(&model.DocumentRenderLog{}).Where("document_id = ?", documentID)

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []model.DocumentRenderLog
	if err := q.Order("created_at DESC").
		Offset(pagination.Offset(page.Page, page.Limit)).
		Limit(page.Limit).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	out := make([]logEntity.RenderLog, len(rows))
	for i := range rows {
		out[i] = model.ToEntity(&rows[i])
	}
	return out, total, nil
}
