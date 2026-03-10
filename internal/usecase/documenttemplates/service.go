package documenttemplates

import (
	"context"
	"errors"
	"strings"
	"time"

	entityTmpl "go-document-generator/internal/entity/documenttemplates"
	entityVer "go-document-generator/internal/entity/documenttemplateversions"
	repoTmpl "go-document-generator/internal/repository/documenttemplates"
	repoVer "go-document-generator/internal/repository/documenttemplateversions"
)

type TemplatesService interface {
	// CreateOrAddVersion:
	// - Jika template (berdasarkan code) belum ada: buat template (IsActive=true) + versi 1 (IsPublished=true).
	// - Jika sudah ada: tambahkan versi baru (auto increment), default IsPublished=false.
	CreateOrAddVersion(ctx context.Context, tmpl entityTmpl.DocumentTemplate, ver entityVer.DocumentTemplateVersion) (entityTmpl.DocumentTemplate, entityVer.DocumentTemplateVersion, error)

	// CRUD Templates
	Create(ctx context.Context, tmpl entityTmpl.DocumentTemplate) (entityTmpl.DocumentTemplate, error)
	GetByID(ctx context.Context, id int64) (entityTmpl.DocumentTemplate, error)
	GetByCode(ctx context.Context, code string) (entityTmpl.DocumentTemplate, error)
	List(ctx context.Context) ([]entityTmpl.DocumentTemplate, error)
	Update(ctx context.Context, tmpl entityTmpl.DocumentTemplate) (entityTmpl.DocumentTemplate, error)
	Delete(ctx context.Context, id int64) error
}

type templatesService struct {
	tmplRepo repoTmpl.DocumentTemplatesRepository
	verRepo  repoVer.DocumentTemplateVersionsRepository
}

func NewTemplatesService(tmplRepo repoTmpl.DocumentTemplatesRepository, verRepo repoVer.DocumentTemplateVersionsRepository) TemplatesService {
	return &templatesService{
		tmplRepo: tmplRepo,
		verRepo:  verRepo,
	}
}

func (s *templatesService) CreateOrAddVersion(ctx context.Context, tmpl entityTmpl.DocumentTemplate, ver entityVer.DocumentTemplateVersion) (entityTmpl.DocumentTemplate, entityVer.DocumentTemplateVersion, error) {
	// Normalisasi input sederhana
	tmpl.Code = strings.TrimSpace(tmpl.Code)
	tmpl.Name = strings.TrimSpace(tmpl.Name)
	tmpl.Engine = strings.TrimSpace(tmpl.Engine)
	tmpl.OutputFormat = strings.TrimSpace(tmpl.OutputFormat)

	if tmpl.Code == "" || tmpl.Name == "" || tmpl.Engine == "" || tmpl.OutputFormat == "" {
		return entityTmpl.DocumentTemplate{}, entityVer.DocumentTemplateVersion{}, errors.New("code, name, engine, output_format wajib diisi")
	}

	// Cek apakah template dengan code sudah ada
	existing, err := s.tmplRepo.GetByCode(ctx, tmpl.Code)
	if err != nil {
		// Asumsikan err berarti belum ada (repo mengembalikan "not found")
		now := time.Now()
		if tmpl.CreatedAt.IsZero() {
			tmpl.CreatedAt = now
		}
		if tmpl.UpdatedAt.IsZero() {
			tmpl.UpdatedAt = now
		}
		tmpl.IsActive = true
		createdTmpl, errCreate := s.tmplRepo.Create(ctx, tmpl)
		if errCreate != nil {
			return entityTmpl.DocumentTemplate{}, entityVer.DocumentTemplateVersion{}, errCreate
		}
		// Siapkan versi 1
		if ver.Version == 0 {
			ver.Version = 1
		}
		ver.TemplateID = createdTmpl.ID
		if !ver.IsPublished {
			ver.IsPublished = true
		}
		createdVer, errVer := s.verRepo.Create(ctx, ver)
		if errVer != nil {
			return entityTmpl.DocumentTemplate{}, entityVer.DocumentTemplateVersion{}, errVer
		}
		return createdTmpl, createdVer, nil
	}

	// Sudah ada template → tambahkan versi baru
	maxv, err := s.verRepo.GetLatestVersionNumber(ctx, existing.ID)
	if err != nil {
		return entityTmpl.DocumentTemplate{}, entityVer.DocumentTemplateVersion{}, err
	}
	if ver.Version == 0 {
		ver.Version = maxv + 1
	}
	ver.TemplateID = existing.ID
	// default tidak dipublish
	createdVer, err := s.verRepo.Create(ctx, ver)
	if err != nil {
		return entityTmpl.DocumentTemplate{}, entityVer.DocumentTemplateVersion{}, err
	}
	return existing, createdVer, nil
}

// CRUD Templates
func (s *templatesService) Create(ctx context.Context, tmpl entityTmpl.DocumentTemplate) (entityTmpl.DocumentTemplate, error) {
	return s.tmplRepo.Create(ctx, tmpl)
}

func (s *templatesService) GetByID(ctx context.Context, id int64) (entityTmpl.DocumentTemplate, error) {
	return s.tmplRepo.GetByID(ctx, id)
}

func (s *templatesService) GetByCode(ctx context.Context, code string) (entityTmpl.DocumentTemplate, error) {
	return s.tmplRepo.GetByCode(ctx, code)
}

func (s *templatesService) List(ctx context.Context) ([]entityTmpl.DocumentTemplate, error) {
	return s.tmplRepo.List(ctx)
}

func (s *templatesService) Update(ctx context.Context, tmpl entityTmpl.DocumentTemplate) (entityTmpl.DocumentTemplate, error) {
	return s.tmplRepo.Update(ctx, tmpl)
}

func (s *templatesService) Delete(ctx context.Context, id int64) error {
	return s.tmplRepo.Delete(ctx, id)
}

