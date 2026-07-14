package template

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	tplEntity "go-document-generator/internal/entity/documenttemplates"
	verEntity "go-document-generator/internal/entity/documenttemplateversions"
	tplrepo "go-document-generator/internal/repository/documenttemplates"
	verrepo "go-document-generator/internal/repository/documenttemplateversions"

	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const defaultTTL = 5 * time.Minute

// CachedTemplateRepo membungkus DocumentTemplatesRepository dengan Redis cache.
type CachedTemplateRepo struct {
	inner  tplrepo.DocumentTemplatesRepository
	redis  *goredis.Client
	ttl    time.Duration
}

func NewCachedTemplateRepo(inner tplrepo.DocumentTemplatesRepository, redis *goredis.Client) tplrepo.DocumentTemplatesRepository {
	return &CachedTemplateRepo{inner: inner, redis: redis, ttl: defaultTTL}
}

func (r *CachedTemplateRepo) GetByCode(ctx context.Context, tx *gorm.DB, code string, tenantID *string) (tplEntity.Template, error) {
	key := templateByCodeKey(code, tenantID)
	if cached, err := getJSON[tplEntity.Template](ctx, r.redis, key); err == nil {
		return cached, nil
	}
	tpl, err := r.inner.GetByCode(ctx, tx, code, tenantID)
	if err != nil {
		return tplEntity.Template{}, err
	}
	_ = setJSON(ctx, r.redis, key, tpl, r.ttl)
	return tpl, nil
}

func (r *CachedTemplateRepo) GetByID(ctx context.Context, tx *gorm.DB, id int64, tenantID *string) (tplEntity.Template, error) {
	key := templateByIDKey(id, tenantID)
	if cached, err := getJSON[tplEntity.Template](ctx, r.redis, key); err == nil {
		return cached, nil
	}
	tpl, err := r.inner.GetByID(ctx, tx, id, tenantID)
	if err != nil {
		return tplEntity.Template{}, err
	}
	_ = setJSON(ctx, r.redis, key, tpl, r.ttl)
	return tpl, nil
}

func (r *CachedTemplateRepo) Create(ctx context.Context, tx *gorm.DB, t tplEntity.Template) (tplEntity.Template, error) {
	return r.inner.Create(ctx, tx, t)
}

func (r *CachedTemplateRepo) List(ctx context.Context, tx *gorm.DB, f tplrepo.ListFilter) ([]tplEntity.Template, int64, error) {
	return r.inner.List(ctx, tx, f)
}

func (r *CachedTemplateRepo) Update(ctx context.Context, tx *gorm.DB, t tplEntity.Template) (tplEntity.Template, error) {
	result, err := r.inner.Update(ctx, tx, t)
	if err == nil {
		_ = r.redis.Del(ctx, templateByIDKey(t.ID, t.TenantID), templateByCodeKey(t.Code, t.TenantID))
	}
	return result, err
}

func (r *CachedTemplateRepo) Deactivate(ctx context.Context, tx *gorm.DB, id int64, tenantID *string, updatedBy *string) error {
	err := r.inner.Deactivate(ctx, tx, id, tenantID, updatedBy)
	if err == nil {
		_ = r.redis.Del(ctx, templateByIDKey(id, tenantID))
	}
	return err
}

// CachedVersionRepo membungkus DocumentTemplateVersionsRepository dengan Redis cache.
type CachedVersionRepo struct {
	inner verrepo.DocumentTemplateVersionsRepository
	redis *goredis.Client
	ttl   time.Duration
}

func NewCachedVersionRepo(inner verrepo.DocumentTemplateVersionsRepository, redis *goredis.Client) verrepo.DocumentTemplateVersionsRepository {
	return &CachedVersionRepo{inner: inner, redis: redis, ttl: defaultTTL}
}

func (r *CachedVersionRepo) GetLatestPublished(ctx context.Context, tx *gorm.DB, templateID int64, tenantID *string) (verEntity.TemplateVersion, error) {
	key := versionLatestKey(templateID, tenantID)
	if cached, err := getJSON[verEntity.TemplateVersion](ctx, r.redis, key); err == nil {
		return cached, nil
	}
	ver, err := r.inner.GetLatestPublished(ctx, tx, templateID, tenantID)
	if err != nil {
		return verEntity.TemplateVersion{}, err
	}
	_ = setJSON(ctx, r.redis, key, ver, r.ttl)
	return ver, nil
}

func (r *CachedVersionRepo) GetByTemplateAndVersion(ctx context.Context, tx *gorm.DB, templateID int64, version int, tenantID *string) (verEntity.TemplateVersion, error) {
	key := versionByNumKey(templateID, version, tenantID)
	if cached, err := getJSON[verEntity.TemplateVersion](ctx, r.redis, key); err == nil {
		return cached, nil
	}
	ver, err := r.inner.GetByTemplateAndVersion(ctx, tx, templateID, version, tenantID)
	if err != nil {
		return verEntity.TemplateVersion{}, err
	}
	_ = setJSON(ctx, r.redis, key, ver, r.ttl)
	return ver, nil
}

func (r *CachedVersionRepo) GetByID(ctx context.Context, tx *gorm.DB, templateID, versionID int64, tenantID *string) (verEntity.TemplateVersion, error) {
	key := versionByIDKey(templateID, versionID, tenantID)
	if cached, err := getJSON[verEntity.TemplateVersion](ctx, r.redis, key); err == nil {
		return cached, nil
	}
	ver, err := r.inner.GetByID(ctx, tx, templateID, versionID, tenantID)
	if err != nil {
		return verEntity.TemplateVersion{}, err
	}
	_ = setJSON(ctx, r.redis, key, ver, r.ttl)
	return ver, nil
}

func (r *CachedVersionRepo) Create(ctx context.Context, tx *gorm.DB, v verEntity.TemplateVersion) (verEntity.TemplateVersion, error) {
	return r.inner.Create(ctx, tx, v)
}

func (r *CachedVersionRepo) ListByTemplateID(ctx context.Context, tx *gorm.DB, templateID int64, tenantID *string, isPublished *bool) ([]verEntity.TemplateVersion, error) {
	return r.inner.ListByTemplateID(ctx, tx, templateID, tenantID, isPublished)
}

func (r *CachedVersionRepo) NextVersionNumber(ctx context.Context, tx *gorm.DB, templateID int64) (int, error) {
	return r.inner.NextVersionNumber(ctx, tx, templateID)
}

func (r *CachedVersionRepo) UnpublishOthers(ctx context.Context, tx *gorm.DB, templateID, exceptVersionID int64) error {
	return r.inner.UnpublishOthers(ctx, tx, templateID, exceptVersionID)
}

func (r *CachedVersionRepo) Publish(ctx context.Context, tx *gorm.DB, templateID, versionID int64, tenantID *string) (verEntity.TemplateVersion, error) {
	result, err := r.inner.Publish(ctx, tx, templateID, versionID, tenantID)
	if err == nil {
		_ = r.redis.Del(ctx, versionLatestKey(templateID, tenantID))
	}
	return result, err
}


// helpers

func tenantStr(tenantID *string) string {
	if tenantID == nil {
		return "global"
	}
	return *tenantID
}

func templateByCodeKey(code string, tenantID *string) string {
	return fmt.Sprintf("tpl:code:%s:%s", tenantStr(tenantID), code)
}

func templateByIDKey(id int64, tenantID *string) string {
	return fmt.Sprintf("tpl:id:%s:%d", tenantStr(tenantID), id)
}

func versionLatestKey(templateID int64, tenantID *string) string {
	return fmt.Sprintf("ver:latest:%s:%d", tenantStr(tenantID), templateID)
}

func versionByNumKey(templateID int64, version int, tenantID *string) string {
	return fmt.Sprintf("ver:num:%s:%d:%d", tenantStr(tenantID), templateID, version)
}

func versionByIDKey(templateID, versionID int64, tenantID *string) string {
	return fmt.Sprintf("ver:id:%s:%d:%d", tenantStr(tenantID), templateID, versionID)
}

func getJSON[T any](ctx context.Context, r *goredis.Client, key string) (T, error) {
	var zero T
	val, err := r.Get(ctx, key).Bytes()
	if err != nil {
		return zero, err
	}
	var out T
	if err := json.Unmarshal(val, &out); err != nil {
		return zero, err
	}
	return out, nil
}

func setJSON(ctx context.Context, r *goredis.Client, key string, val any, ttl time.Duration) error {
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return r.Set(ctx, key, b, ttl).Err()
}
