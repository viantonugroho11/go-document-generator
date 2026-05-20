package postgres

import (
	"context"
	"time"

	cbEntity "go-document-generator/internal/entity/documentcallbackattempts"
	repo "go-document-generator/internal/repository/documentcallbackattempts"
	"go-document-generator/internal/repository/documentcallbackattempts/model"
	"go-document-generator/internal/shared/pagination"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewDocumentCallbackAttemptsRepository(db *gorm.DB) repo.DocumentCallbackAttemptsRepository {
	return &repository{db: db}
}

func (r *repository) conn(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *repository) Create(ctx context.Context, tx *gorm.DB, a cbEntity.CallbackAttempt) (cbEntity.CallbackAttempt, error) {
	m := model.ToModel(a)
	if m.AttemptedAt.IsZero() {
		m.AttemptedAt = time.Now().UTC()
	}
	if err := r.conn(tx).WithContext(ctx).Create(&m).Error; err != nil {
		return cbEntity.CallbackAttempt{}, err
	}
	return model.ToEntity(&m), nil
}

func (r *repository) ListByDocumentID(ctx context.Context, tx *gorm.DB, documentID int64, page pagination.Params) ([]cbEntity.CallbackAttempt, int64, error) {
	q := r.conn(tx).WithContext(ctx).Model(&model.DocumentCallbackAttempt{}).Where("document_id = ?", documentID)

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []model.DocumentCallbackAttempt
	if err := q.Order("attempted_at DESC").
		Offset(pagination.Offset(page.Page, page.Limit)).
		Limit(page.Limit).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	out := make([]cbEntity.CallbackAttempt, len(rows))
	for i := range rows {
		out[i] = model.ToEntity(&rows[i])
	}
	return out, total, nil
}
