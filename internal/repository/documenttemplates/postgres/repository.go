package postgres

import (
	"context"
	"errors"
	"strings"
	"time"

	tplEntity "go-document-generator/internal/entity/documenttemplates"
	repo "go-document-generator/internal/repository/documenttemplates"
	"go-document-generator/internal/repository/documenttemplates/model"
	"go-document-generator/internal/shared/apperror"
	"go-document-generator/internal/shared/pagination"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewDocumentTemplatesRepository(db *gorm.DB) repo.DocumentTemplatesRepository {
	return &repository{db: db}
}

func (r *repository) conn(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *repository) Create(ctx context.Context, tx *gorm.DB, t tplEntity.Template) (tplEntity.Template, error) {
	m := model.ToModel(t)
	m.CreatedAt = time.Now().UTC()
	m.UpdatedAt = m.CreatedAt
	if err := r.conn(tx).WithContext(ctx).Create(&m).Error; err != nil {
		return tplEntity.Template{}, err
	}
	return model.ToEntity(&m), nil
}

func (r *repository) GetByID(ctx context.Context, tx *gorm.DB, id int64, tenantID *string) (tplEntity.Template, error) {
	var m model.DocumentTemplate
	q := r.conn(tx).WithContext(ctx).Where("id = ?", id)
	if tenantID != nil {
		q = q.Where("tenant_id = ?", *tenantID)
	}
	if err := q.First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tplEntity.Template{}, apperror.ErrNotFound
		}
		return tplEntity.Template{}, err
	}
	return model.ToEntity(&m), nil
}

func (r *repository) GetByCode(ctx context.Context, tx *gorm.DB, code string, tenantID *string) (tplEntity.Template, error) {
	var m model.DocumentTemplate
	q := r.conn(tx).WithContext(ctx).Where("code = ?", code)
	if tenantID != nil {
		q = q.Where("tenant_id = ?", *tenantID)
	} else {
		q = q.Where("tenant_id IS NULL")
	}
	if err := q.First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tplEntity.Template{}, apperror.ErrNotFound
		}
		return tplEntity.Template{}, err
	}
	return model.ToEntity(&m), nil
}

func (r *repository) List(ctx context.Context, tx *gorm.DB, f repo.ListFilter) ([]tplEntity.Template, int64, error) {
	q := r.conn(tx).WithContext(ctx).Model(&model.DocumentTemplate{})
	if f.TenantID != nil {
		q = q.Where("tenant_id = ?", *f.TenantID)
	}
	if f.Code != "" {
		q = q.Where("code = ?", f.Code)
	}
	if f.Category != "" {
		q = q.Where("category = ?", f.Category)
	}
	if f.IsActive != nil {
		q = q.Where("is_active = ?", *f.IsActive)
	}
	if f.Engine != "" {
		q = q.Where("engine = ?", f.Engine)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sort := "created_at DESC"
	if s := strings.TrimSpace(f.Page.Sort); s != "" {
		sort = normalizeSort(s)
	}
	q = q.Order(sort).Offset(pagination.Offset(f.Page.Page, f.Page.Limit)).Limit(f.Page.Limit)

	var rows []model.DocumentTemplate
	if err := q.Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	out := make([]tplEntity.Template, len(rows))
	for i := range rows {
		out[i] = model.ToEntity(&rows[i])
	}
	return out, total, nil
}

func (r *repository) Update(ctx context.Context, tx *gorm.DB, t tplEntity.Template) (tplEntity.Template, error) {
	updates := map[string]any{
		"updated_at": time.Now().UTC(),
	}
	if t.Name != "" {
		updates["name"] = t.Name
	}
	if t.Description != nil {
		updates["description"] = t.Description
	}
	if t.Engine != "" {
		updates["engine"] = t.Engine
	}
	if t.DefaultFormat != "" {
		updates["default_format"] = t.DefaultFormat
	}
	if t.Category != nil {
		updates["category"] = t.Category
	}
	updates["is_active"] = t.IsActive
	if t.UpdatedBy != nil {
		updates["updated_by"] = t.UpdatedBy
	}

	q := r.conn(tx).WithContext(ctx).Model(&model.DocumentTemplate{}).Where("id = ?", t.ID)
	if t.TenantID != nil {
		q = q.Where("tenant_id = ?", *t.TenantID)
	}
	res := q.Updates(updates)
	if res.Error != nil {
		return tplEntity.Template{}, res.Error
	}
	if res.RowsAffected == 0 {
		return tplEntity.Template{}, apperror.ErrNotFound
	}
	return r.GetByID(ctx, tx, t.ID, t.TenantID)
}

func (r *repository) Deactivate(ctx context.Context, tx *gorm.DB, id int64, tenantID *string, updatedBy *string) error {
	updates := map[string]any{
		"is_active":  false,
		"updated_at": time.Now().UTC(),
	}
	if updatedBy != nil {
		updates["updated_by"] = updatedBy
	}
	q := r.conn(tx).WithContext(ctx).Model(&model.DocumentTemplate{}).Where("id = ?", id)
	if tenantID != nil {
		q = q.Where("tenant_id = ?", *tenantID)
	}
	res := q.Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return apperror.ErrNotFound
	}
	return nil
}

func normalizeSort(sort string) string {
	if strings.HasPrefix(sort, "-") {
		return strings.TrimPrefix(sort, "-") + " DESC"
	}
	return sort + " ASC"
}
