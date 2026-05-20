package documenttemplateversions

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"strings"

	verEntity "go-document-generator/internal/entity/documenttemplateversions"
	begin "go-document-generator/internal/repository/begin"
	tplrepo "go-document-generator/internal/repository/documenttemplates"
	verrepo "go-document-generator/internal/repository/documenttemplateversions"
	"go-document-generator/internal/shared/apperror"
)

type Service interface {
	Create(ctx context.Context, templateID int64, tenantID *string, v verEntity.TemplateVersion) (verEntity.TemplateVersion, error)
	GetByID(ctx context.Context, templateID, versionID int64, tenantID *string) (verEntity.TemplateVersion, error)
	List(ctx context.Context, templateID int64, tenantID *string, isPublished *bool) ([]verEntity.TemplateVersion, error)
	Publish(ctx context.Context, templateID, versionID int64, tenantID *string) (verEntity.TemplateVersion, error)
}

type service struct {
	versions  verrepo.DocumentTemplateVersionsRepository
	templates tplrepo.DocumentTemplatesRepository
	txManager begin.BeginRepository
	publisher VersionEventPublisher
}

func NewService(
	versions verrepo.DocumentTemplateVersionsRepository,
	templates tplrepo.DocumentTemplatesRepository,
	tx begin.BeginRepository,
	publisher VersionEventPublisher,
) Service {
	if publisher == nil {
		publisher = NoopVersionPublisher()
	}
	return &service{versions: versions, templates: templates, txManager: tx, publisher: publisher}
}

func (s *service) Create(ctx context.Context, templateID int64, tenantID *string, v verEntity.TemplateVersion) (verEntity.TemplateVersion, error) {
	if templateID <= 0 {
		return verEntity.TemplateVersion{}, apperror.ErrInvalidInput
	}
	if strings.TrimSpace(v.Content) == "" {
		return verEntity.TemplateVersion{}, errors.New("content is required")
	}
	if v.OutputFormat == "" {
		return verEntity.TemplateVersion{}, errors.New("output_format is required")
	}

	if _, err := s.templates.GetByID(ctx, nil, templateID, tenantID); err != nil {
		return verEntity.TemplateVersion{}, mapRepoErr(err)
	}

	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return verEntity.TemplateVersion{}, err
	}
	defer func() {
		if err != nil {
			_ = s.txManager.Rollback(ctx, tx)
		}
	}()

	next, err := s.versions.NextVersionNumber(ctx, tx, templateID)
	if err != nil {
		return verEntity.TemplateVersion{}, err
	}
	v.TemplateID = templateID
	v.TenantID = tenantID
	v.Version = next
	sum := sha256.Sum256([]byte(v.Content))
	chk := hex.EncodeToString(sum[:])
	v.Checksum = &chk

	created, err := s.versions.Create(ctx, tx, v)
	if err != nil {
		return verEntity.TemplateVersion{}, err
	}
	if err = s.txManager.Commit(ctx, tx); err != nil {
		return verEntity.TemplateVersion{}, err
	}
	if pubErr := s.publisher.PublishVersionCreated(ctx, created); pubErr != nil {
		log.Printf("documenttemplateversions: PublishVersionCreated: %v", pubErr)
	}
	return created, nil
}

func (s *service) GetByID(ctx context.Context, templateID, versionID int64, tenantID *string) (verEntity.TemplateVersion, error) {
	v, err := s.versions.GetByID(ctx, nil, templateID, versionID, tenantID)
	return v, mapRepoErr(err)
}

func (s *service) List(ctx context.Context, templateID int64, tenantID *string, isPublished *bool) ([]verEntity.TemplateVersion, error) {
	if _, err := s.templates.GetByID(ctx, nil, templateID, tenantID); err != nil {
		return nil, mapRepoErr(err)
	}
	return s.versions.ListByTemplateID(ctx, nil, templateID, tenantID, isPublished)
}

func (s *service) Publish(ctx context.Context, templateID, versionID int64, tenantID *string) (verEntity.TemplateVersion, error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return verEntity.TemplateVersion{}, err
	}
	defer func() {
		if err != nil {
			_ = s.txManager.Rollback(ctx, tx)
		}
	}()

	if err = s.versions.UnpublishOthers(ctx, tx, templateID, versionID); err != nil {
		return verEntity.TemplateVersion{}, err
	}
	published, err := s.versions.Publish(ctx, tx, templateID, versionID, tenantID)
	if err != nil {
		return verEntity.TemplateVersion{}, mapRepoErr(err)
	}
	if err = s.txManager.Commit(ctx, tx); err != nil {
		return verEntity.TemplateVersion{}, err
	}
	if pubErr := s.publisher.PublishVersionPublished(ctx, published); pubErr != nil {
		log.Printf("documenttemplateversions: PublishVersionPublished: %v", pubErr)
	}
	return published, nil
}

func mapRepoErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, apperror.ErrNotFound) {
		return apperror.ErrNotFound
	}
	return err
}
