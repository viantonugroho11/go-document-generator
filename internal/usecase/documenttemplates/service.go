package documenttemplates

import (
	"context"
	"errors"
	"log"
	"strings"

	tplEntity "go-document-generator/internal/entity/documenttemplates"
	"go-document-generator/internal/entity/enums"
	begin "go-document-generator/internal/repository/begin"
	repo "go-document-generator/internal/repository/documenttemplates"
	"go-document-generator/internal/shared/apperror"
	"go-document-generator/internal/shared/pagination"
)

type Service interface {
	Create(ctx context.Context, t tplEntity.Template) (tplEntity.Template, error)
	GetByID(ctx context.Context, id int64, tenantID *string) (tplEntity.Template, error)
	List(ctx context.Context, f repo.ListFilter) ([]tplEntity.Template, pagination.Meta, error)
	Patch(ctx context.Context, t tplEntity.Template) (tplEntity.Template, error)
	Deactivate(ctx context.Context, id int64, tenantID *string, updatedBy *string) error
}

type service struct {
	repo      repo.DocumentTemplatesRepository
	txManager begin.BeginRepository
	publisher TemplateEventPublisher
}

func NewService(repo repo.DocumentTemplatesRepository, tx begin.BeginRepository, publisher TemplateEventPublisher) Service {
	if publisher == nil {
		publisher = NoopTemplatePublisher()
	}
	return &service{repo: repo, txManager: tx, publisher: publisher}
}

func (s *service) Create(ctx context.Context, t tplEntity.Template) (tplEntity.Template, error) {
	if err := validateTemplate(t, true); err != nil {
		return tplEntity.Template{}, err
	}
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return tplEntity.Template{}, err
	}
	defer func() {
		if err != nil {
			_ = s.txManager.Rollback(ctx, tx)
		}
	}()

	created, err := s.repo.Create(ctx, tx, t)
	if err != nil {
		return tplEntity.Template{}, mapRepoErr(err)
	}
	if err = s.txManager.Commit(ctx, tx); err != nil {
		return tplEntity.Template{}, err
	}

	if pubErr := s.publisher.PublishTemplateCreated(ctx, created); pubErr != nil {
		log.Printf("documenttemplates: PublishTemplateCreated: %v", pubErr)
	}
	return created, nil
}

func (s *service) GetByID(ctx context.Context, id int64, tenantID *string) (tplEntity.Template, error) {
	if id <= 0 {
		return tplEntity.Template{}, apperror.ErrInvalidInput
	}
	t, err := s.repo.GetByID(ctx, nil, id, tenantID)
	return t, mapRepoErr(err)
}

func (s *service) List(ctx context.Context, f repo.ListFilter) ([]tplEntity.Template, pagination.Meta, error) {
	f.Page = pagination.Normalize(f.Page.Page, f.Page.Limit)
	items, total, err := s.repo.List(ctx, nil, f)
	if err != nil {
		return nil, pagination.Meta{}, err
	}
	return items, pagination.Meta{Page: f.Page.Page, Limit: f.Page.Limit, Total: total}, nil
}

func (s *service) Patch(ctx context.Context, t tplEntity.Template) (tplEntity.Template, error) {
	if t.ID <= 0 {
		return tplEntity.Template{}, apperror.ErrInvalidInput
	}
	updated, err := s.repo.Update(ctx, nil, t)
	if err != nil {
		return tplEntity.Template{}, mapRepoErr(err)
	}
	if pubErr := s.publisher.PublishTemplateUpdated(ctx, updated); pubErr != nil {
		log.Printf("documenttemplates: PublishTemplateUpdated: %v", pubErr)
	}
	return updated, nil
}

func (s *service) Deactivate(ctx context.Context, id int64, tenantID *string, updatedBy *string) error {
	if err := s.repo.Deactivate(ctx, nil, id, tenantID, updatedBy); err != nil {
		return mapRepoErr(err)
	}
	return nil
}

func validateTemplate(t tplEntity.Template, creating bool) error {
	if strings.TrimSpace(t.Code) == "" {
		return errors.New("code is required")
	}
	if strings.TrimSpace(t.Name) == "" {
		return errors.New("name is required")
	}
	if creating {
		switch t.Engine {
		case enums.TemplateEngineHandlebars, enums.TemplateEngineMustache, enums.TemplateEngineHTML:
		default:
			return errors.New("invalid engine")
		}
		switch t.DefaultFormat {
		case enums.OutputFormatPDF, enums.OutputFormatHTML, enums.OutputFormatDOCX:
		default:
			return errors.New("invalid default_format")
		}
	}
	return nil
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
