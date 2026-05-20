package postgres

import (
	"context"
	"errors"
	"strings"
	"time"

	docEntity "go-document-generator/internal/entity/documents"
	repo "go-document-generator/internal/repository/documents"
	"go-document-generator/internal/repository/documents/model"
	"go-document-generator/internal/shared/apperror"
	"go-document-generator/internal/shared/pagination"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewDocumentsRepository(db *gorm.DB) repo.DocumentsRepository {
	return &repository{db: db}
}

func (r *repository) conn(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *repository) Create(ctx context.Context, tx *gorm.DB, d docEntity.Document) (docEntity.Document, error) {
	m := model.ToModel(d)
	now := time.Now().UTC()
	m.CreatedAt = now
	m.UpdatedAt = now
	if err := r.conn(tx).WithContext(ctx).Create(&m).Error; err != nil {
		return docEntity.Document{}, err
	}
	return model.ToEntity(&m), nil
}

func (r *repository) GetByID(ctx context.Context, tx *gorm.DB, id int64, tenantID *string) (docEntity.Document, error) {
	var m model.Document
	q := r.conn(tx).WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id)
	if tenantID != nil {
		q = q.Where("tenant_id = ?", *tenantID)
	}
	if err := q.First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return docEntity.Document{}, apperror.ErrNotFound
		}
		return docEntity.Document{}, err
	}
	return model.ToEntity(&m), nil
}

func (r *repository) GetByRequestID(ctx context.Context, tx *gorm.DB, requestID string, tenantID *string) (docEntity.Document, error) {
	var m model.Document
	q := r.conn(tx).WithContext(ctx).Where("request_id = ? AND deleted_at IS NULL", requestID)
	if tenantID != nil {
		q = q.Where("tenant_id = ?", *tenantID)
	} else {
		q = q.Where("tenant_id IS NULL")
	}
	if err := q.First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return docEntity.Document{}, apperror.ErrNotFound
		}
		return docEntity.Document{}, err
	}
	return model.ToEntity(&m), nil
}

func (r *repository) List(ctx context.Context, tx *gorm.DB, f repo.ListFilter) ([]docEntity.Document, int64, error) {
	q := r.conn(tx).WithContext(ctx).Model(&model.Document{}).Where("deleted_at IS NULL")
	if f.TenantID != nil {
		q = q.Where("tenant_id = ?", *f.TenantID)
	}
	if f.RequestID != "" {
		q = q.Where("request_id = ?", f.RequestID)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.TemplateCode != "" {
		q = q.Where("template_code = ?", f.TemplateCode)
	}
	if f.DmsStatus != "" {
		q = q.Where("dms_status = ?", f.DmsStatus)
	}
	if f.CallbackStatus != "" {
		q = q.Where("callback_status = ?", f.CallbackStatus)
	}
	if f.CreatedFrom != nil {
		q = q.Where("created_at >= ?", *f.CreatedFrom)
	}
	if f.CreatedTo != nil {
		q = q.Where("created_at <= ?", *f.CreatedTo)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sort := "created_at DESC"
	if s := strings.TrimSpace(f.Page.Sort); s != "" {
		if strings.HasPrefix(s, "-") {
			sort = strings.TrimPrefix(s, "-") + " DESC"
		} else {
			sort = s + " ASC"
		}
	}
	q = q.Order(sort).Offset(pagination.Offset(f.Page.Page, f.Page.Limit)).Limit(f.Page.Limit)

	var rows []model.Document
	if err := q.Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	out := make([]docEntity.Document, len(rows))
	for i := range rows {
		out[i] = model.ToEntity(&rows[i])
	}
	return out, total, nil
}

func (r *repository) Update(ctx context.Context, tx *gorm.DB, d docEntity.Document) (docEntity.Document, error) {
	m := model.ToModel(d)
	m.UpdatedAt = time.Now().UTC()
	if err := r.conn(tx).WithContext(ctx).Save(&m).Error; err != nil {
		return docEntity.Document{}, err
	}
	return model.ToEntity(&m), nil
}

func (r *repository) SoftDelete(ctx context.Context, tx *gorm.DB, id int64, tenantID *string) error {
	now := time.Now().UTC()
	q := r.conn(tx).WithContext(ctx).
		Model(&model.Document{}).
		Where("id = ? AND deleted_at IS NULL", id)
	if tenantID != nil {
		q = q.Where("tenant_id = ?", *tenantID)
	}
	res := q.Update("deleted_at", now)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return apperror.ErrNotFound
	}
	return nil
}
