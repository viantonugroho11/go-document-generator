package documents

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	docEntity "go-document-generator/internal/entity/documents"
	verEntity "go-document-generator/internal/entity/documenttemplateversions"
	"go-document-generator/internal/entity/enums"
	begin "go-document-generator/internal/repository/begin"
	docrepo "go-document-generator/internal/repository/documents"
	tplrepo "go-document-generator/internal/repository/documenttemplates"
	verrepo "go-document-generator/internal/repository/documenttemplateversions"
	"go-document-generator/internal/shared/apperror"
	"go-document-generator/internal/shared/pagination"
	"go-document-generator/internal/shared/validators"
	"go-document-generator/internal/usecase/documents/states"
	"go-document-generator/internal/usecase/documents/transitions"
)

type CreateInput struct {
	TenantID        *string
	RequestID       string
	TemplateCode    string
	TemplateVersion *int
	OutputFormat    enums.OutputFormat
	Payload         map[string]any
	Metadata        map[string]any
	StoreToDms      bool
	HasCallback     bool
	CallbackURL     *string
	ExpiredAt       *time.Time
	CreatedBy       *string
}

type Service interface {
	Create(ctx context.Context, in CreateInput) (docEntity.Document, bool, error)
	GetByID(ctx context.Context, id int64, tenantID *string) (docEntity.Document, error)
	Patch(ctx context.Context, d docEntity.Document) (docEntity.Document, error)
	GetByRequestID(ctx context.Context, requestID string, tenantID *string) (docEntity.Document, error)
	List(ctx context.Context, f docrepo.ListFilter) ([]docEntity.Document, pagination.Meta, error)
	Cancel(ctx context.Context, id int64, tenantID *string) (docEntity.Document, error)
	Retry(ctx context.Context, id int64, tenantID *string) (docEntity.Document, error)
	SoftDelete(ctx context.Context, id int64, tenantID *string) error
	DownloadURL(ctx context.Context, id int64, tenantID *string) (string, error)
}

type service struct {
	docs          docrepo.DocumentsRepository
	templates     tplrepo.DocumentTemplatesRepository
	versions      verrepo.DocumentTemplateVersionsRepository
	txManager     begin.BeginRepository
	publisher     DocumentEventPublisher
	selector      GeneratorSelector
	stateMachine  states.IDocumentStateMachineFactory
}

func NewService(
	docs docrepo.DocumentsRepository,
	templates tplrepo.DocumentTemplatesRepository,
	versions verrepo.DocumentTemplateVersionsRepository,
	tx begin.BeginRepository,
	publisher DocumentEventPublisher,
	selector GeneratorSelector,
) Service {
	if publisher == nil {
		publisher = NoopDocumentPublisher()
	}
	deps := transitions.Deps{Templates: templates, Versions: versions, Selector: adaptSelector(selector)}
	smFactory := states.NewDocumentStateMachineFactory(BuildStateHandlers(deps))

	return &service{
		docs:         docs,
		templates:    templates,
		versions:     versions,
		txManager:    tx,
		publisher:    publisher,
		selector:     selector,
		stateMachine: smFactory,
	}
}

func (s *service) Create(ctx context.Context, in CreateInput) (docEntity.Document, bool, error) {
	if strings.TrimSpace(in.RequestID) == "" || strings.TrimSpace(in.TemplateCode) == "" {
		return docEntity.Document{}, false, apperror.ErrInvalidInput
	}
	if in.Payload == nil {
		return docEntity.Document{}, false, errors.New("payload is required")
	}

	existing, err := s.docs.GetByRequestID(ctx, nil, in.RequestID, in.TenantID)
	if err == nil {
		return existing, true, nil
	}
	if !errors.Is(err, apperror.ErrNotFound) {
		return docEntity.Document{}, false, err
	}

	tpl, err := s.templates.GetByCode(ctx, nil, in.TemplateCode, in.TenantID)
	if err != nil {
		return docEntity.Document{}, false, mapRepoErr(err)
	}
	if !tpl.IsActive {
		return docEntity.Document{}, false, apperror.ErrNotFound
	}

	var ver verEntity.TemplateVersion
	if in.TemplateVersion != nil {
		ver, err = s.versions.GetByTemplateAndVersion(ctx, nil, tpl.ID, *in.TemplateVersion, in.TenantID)
	} else {
		ver, err = s.versions.GetLatestPublished(ctx, nil, tpl.ID, in.TenantID)
	}
	if err != nil {
		return docEntity.Document{}, false, mapRepoErr(err)
	}
	if !ver.IsPublished {
		return docEntity.Document{}, false, apperror.ErrNotFound
	}

	if len(ver.Schema) > 0 {
		if err := validators.ValidateSchema(ver.Schema, in.Payload); err != nil {
			return docEntity.Document{}, false, err
		}
	}

	outFmt := in.OutputFormat
	if outFmt == "" {
		outFmt = ver.OutputFormat
	}
	if outFmt == "" {
		outFmt = tpl.DefaultFormat
	}

	tplID := tpl.ID
	verID := ver.ID
	doc := docEntity.Document{
		TenantID:          in.TenantID,
		RequestID:         in.RequestID,
		TemplateID:        &tplID,
		TemplateVersionID: &verID,
		TemplateCode:      tpl.Code,
		TemplateVersion:   ver.Version,
		Payload:           in.Payload,
		Metadata:          in.Metadata,
		Status:            enums.DocumentStatusQueued,
		OutputFormat:      outFmt,
		StoreToDms:        in.StoreToDms,
		DmsStatus:         enums.DmsStatusNotSent,
		HasCallback:       in.HasCallback,
		CallbackURL:       in.CallbackURL,
		CallbackStatus:    enums.CallbackStatusPending,
		ExpiredAt:         in.ExpiredAt,
		CreatedBy:         in.CreatedBy,
	}

	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return docEntity.Document{}, false, err
	}
	defer func() {
		if err != nil {
			_ = s.txManager.Rollback(ctx, tx)
		}
	}()

	created, err := s.docs.Create(ctx, tx, doc)
	if err != nil {
		return docEntity.Document{}, false, err
	}
	if err = s.txManager.Commit(ctx, tx); err != nil {
		return docEntity.Document{}, false, err
	}
	if pubErr := s.publisher.PublishDocumentQueued(ctx, created); pubErr != nil {
		log.Printf("documents: PublishDocumentQueued: %v", pubErr)
	}
	return created, false, nil
}

func (s *service) GetByID(ctx context.Context, id int64, tenantID *string) (docEntity.Document, error) {
	d, err := s.docs.GetByID(ctx, nil, id, tenantID)
	return d, mapRepoErr(err)
}

func (s *service) Patch(ctx context.Context, patch docEntity.Document) (docEntity.Document, error) {
	if patch.ID <= 0 {
		return docEntity.Document{}, apperror.ErrInvalidInput
	}
	existing, err := s.docs.GetByID(ctx, nil, patch.ID, patch.TenantID)
	if err != nil {
		return docEntity.Document{}, mapRepoErr(err)
	}

	updated := mergeDocumentPatch(existing, patch)
	result, err := s.applyStateMachine(ctx, existing, updated)
	if err != nil {
		return docEntity.Document{}, err
	}
	return s.docs.Update(ctx, nil, result)
}

// mergeDocumentPatch menggabungkan field patch ke dokumen existing.
func mergeDocumentPatch(existing, patch docEntity.Document) docEntity.Document {
	out := existing
	if patch.Payload != nil {
		out.Payload = patch.Payload
	}
	if patch.Metadata != nil {
		out.Metadata = patch.Metadata
	}
	if patch.OutputFormat != "" {
		out.OutputFormat = patch.OutputFormat
	}
	if patch.Status != "" && patch.Status != existing.Status {
		out.Status = patch.Status
	}
	if patch.StoreToDms != existing.StoreToDms {
		out.StoreToDms = patch.StoreToDms
	}
	if patch.HasCallback != existing.HasCallback {
		out.HasCallback = patch.HasCallback
	}
	if patch.CallbackURL != nil {
		out.CallbackURL = patch.CallbackURL
	}
	if patch.ExpiredAt != nil {
		out.ExpiredAt = patch.ExpiredAt
	}
	if patch.ErrorMessage != nil {
		out.ErrorMessage = patch.ErrorMessage
	}
	return out
}

func (s *service) GetByRequestID(ctx context.Context, requestID string, tenantID *string) (docEntity.Document, error) {
	d, err := s.docs.GetByRequestID(ctx, nil, requestID, tenantID)
	return d, mapRepoErr(err)
}

func (s *service) List(ctx context.Context, f docrepo.ListFilter) ([]docEntity.Document, pagination.Meta, error) {
	f.Page = pagination.Normalize(f.Page.Page, f.Page.Limit)
	items, total, err := s.docs.List(ctx, nil, f)
	if err != nil {
		return nil, pagination.Meta{}, err
	}
	return items, pagination.Meta{Page: f.Page.Page, Limit: f.Page.Limit, Total: total}, nil
}

func (s *service) Cancel(ctx context.Context, id int64, tenantID *string) (docEntity.Document, error) {
	existing, err := s.docs.GetByID(ctx, nil, id, tenantID)
	if err != nil {
		return docEntity.Document{}, mapRepoErr(err)
	}
	result, err := s.transitionDocument(ctx, existing, enums.DocumentStatusCancelled)
	if err != nil {
		return docEntity.Document{}, err
	}
	return s.docs.Update(ctx, nil, result)
}

func (s *service) Retry(ctx context.Context, id int64, tenantID *string) (docEntity.Document, error) {
	existing, err := s.docs.GetByID(ctx, nil, id, tenantID)
	if err != nil {
		return docEntity.Document{}, mapRepoErr(err)
	}
	result, err := s.transitionDocument(ctx, existing, enums.DocumentStatusQueued)
	if err != nil {
		return docEntity.Document{}, err
	}
	updated, err := s.docs.Update(ctx, nil, result)
	if err != nil {
		return docEntity.Document{}, err
	}
	if pubErr := s.publisher.PublishDocumentRetried(ctx, updated); pubErr != nil {
		log.Printf("documents: PublishDocumentRetried: %v", pubErr)
	}
	return updated, nil
}

func (s *service) SoftDelete(ctx context.Context, id int64, tenantID *string) error {
	return mapRepoErr(s.docs.SoftDelete(ctx, nil, id, tenantID))
}

func (s *service) DownloadURL(ctx context.Context, id int64, tenantID *string) (string, error) {
	d, err := s.docs.GetByID(ctx, nil, id, tenantID)
	if err != nil {
		return "", mapRepoErr(err)
	}
	if d.Status != enums.DocumentStatusGenerated {
		return "", apperror.ErrInvalidState
	}
	if d.FilePath == nil || strings.TrimSpace(*d.FilePath) == "" {
		return "", apperror.ErrNotFound
	}
	return *d.FilePath, nil
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
