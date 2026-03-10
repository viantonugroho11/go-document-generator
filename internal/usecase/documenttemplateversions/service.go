package documenttemplateversions

import (
	"context"

	entityVer "go-document-generator/internal/entity/documenttemplateversions"
	repoVer "go-document-generator/internal/repository/documenttemplateversions"
)

type VersionsService interface {
	Create(ctx context.Context, ver entityVer.DocumentTemplateVersion) (entityVer.DocumentTemplateVersion, error)
	GetByID(ctx context.Context, id int64) (entityVer.DocumentTemplateVersion, error)
	List(ctx context.Context) ([]entityVer.DocumentTemplateVersion, error)
	ListByTemplateID(ctx context.Context, templateID int64) ([]entityVer.DocumentTemplateVersion, error)
	Update(ctx context.Context, ver entityVer.DocumentTemplateVersion) (entityVer.DocumentTemplateVersion, error)
	Delete(ctx context.Context, id int64) error
}

type versionsService struct {
	repo repoVer.DocumentTemplateVersionsRepository
}

func NewVersionsService(repo repoVer.DocumentTemplateVersionsRepository) VersionsService {
	return &versionsService{repo: repo}
}

func (s *versionsService) Create(ctx context.Context, ver entityVer.DocumentTemplateVersion) (entityVer.DocumentTemplateVersion, error) {
	return s.repo.Create(ctx, ver)
}

func (s *versionsService) GetByID(ctx context.Context, id int64) (entityVer.DocumentTemplateVersion, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *versionsService) List(ctx context.Context) ([]entityVer.DocumentTemplateVersion, error) {
	return s.repo.List(ctx)
}

func (s *versionsService) ListByTemplateID(ctx context.Context, templateID int64) ([]entityVer.DocumentTemplateVersion, error) {
	return s.repo.ListByTemplateID(ctx, templateID)
}

func (s *versionsService) Update(ctx context.Context, ver entityVer.DocumentTemplateVersion) (entityVer.DocumentTemplateVersion, error) {
	return s.repo.Update(ctx, ver)
}

func (s *versionsService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

