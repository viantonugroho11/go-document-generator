package documents

import (
	"context"
	"errors"
	"fmt"
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
	sharedStorage "go-document-generator/internal/shared/storage"
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

type BulkCreateItem struct {
	Input  CreateInput
	Doc    docEntity.Document
	Replay bool
	Err    error
}

type Service interface {
	Create(ctx context.Context, in CreateInput) (docEntity.Document, bool, error)
	BulkCreate(ctx context.Context, inputs []CreateInput) []BulkCreateItem
	GetByID(ctx context.Context, id int64, tenantID *string) (docEntity.Document, error)
	Patch(ctx context.Context, d docEntity.Document) (docEntity.Document, error)
	GetByRequestID(ctx context.Context, requestID string, tenantID *string) (docEntity.Document, error)
	List(ctx context.Context, f docrepo.ListFilter) ([]docEntity.Document, pagination.Meta, error)
	Cancel(ctx context.Context, id int64, tenantID *string) (docEntity.Document, error)
	Retry(ctx context.Context, id int64, tenantID *string) (docEntity.Document, error)
	SoftDelete(ctx context.Context, id int64, tenantID *string) error
	DownloadURL(ctx context.Context, id int64, tenantID *string) (string, error)
	Preview(ctx context.Context, templateID, versionID int64, tenantID *string, payload map[string]any) ([]byte, string, error)
	// ZipDocuments mengambil file dari banyak dokumen, membuat arsip ZIP, mengembalikan URL download.
	ZipDocuments(ctx context.Context, ids []int64, tenantID *string, label string) (string, error)
	// MergeDocuments menggabungkan banyak dokumen format sama menjadi satu file, mengembalikan URL download.
	MergeDocuments(ctx context.Context, ids []int64, tenantID *string, label string) (string, error)
	// Process dipanggil oleh Kafka consumer untuk menjalankan generation pipeline:
	// QUEUED → PROCESSING → GENERATED (atau FAILED bila error).
	Process(ctx context.Context, id int64, tenantID *string) error
}

// StorageProvider abstraksi storage untuk usecase layer.
type StorageProvider interface {
	PresignedURL(ctx context.Context, path string, ttl time.Duration) (string, error)
	ProviderName() enums.StorageProvider
	Download(ctx context.Context, path string) ([]byte, error)
	Zip(ctx context.Context, documentID int64, requestID string, entries []sharedStorage.ZipEntry) (path, fileName string, err error)
	Compose(ctx context.Context, documentID int64, requestID string, srcPaths []string, ext string) (path, fileName string, err error)
}

type service struct {
	docs         docrepo.DocumentsRepository
	templates    tplrepo.DocumentTemplatesRepository
	versions     verrepo.DocumentTemplateVersionsRepository
	txManager    begin.BeginRepository
	publisher    DocumentEventPublisher
	selector     GeneratorSelector
	stateMachine states.IDocumentStateMachineFactory
	storage      StorageProvider
}

func NewService(
	docs docrepo.DocumentsRepository,
	templates tplrepo.DocumentTemplatesRepository,
	versions verrepo.DocumentTemplateVersionsRepository,
	tx begin.BeginRepository,
	publisher DocumentEventPublisher,
	selector GeneratorSelector,
	storageProv StorageProvider,
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
		storage:      storageProv,
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
	if pubErr := s.publisher.PublishDocumentProcess(ctx, created); pubErr != nil {
		log.Printf("documents: PublishDocumentProcess: %v", pubErr)
	}
	return created, false, nil
}

const bulkWorkers = 10

// BulkCreate membuat banyak dokumen secara konkuren (maks bulkWorkers goroutine).
// Tiap item diproses independen; error satu item tidak menghentikan item lain.
func (s *service) BulkCreate(ctx context.Context, inputs []CreateInput) []BulkCreateItem {
	type indexed struct {
		i   int
		res BulkCreateItem
	}
	results := make([]BulkCreateItem, len(inputs))
	out := make(chan indexed, len(inputs))
	sem := make(chan struct{}, bulkWorkers)

	for i, in := range inputs {
		i, in := i, in
		sem <- struct{}{}
		go func() {
			defer func() { <-sem }()
			doc, replay, err := s.Create(ctx, in)
			out <- indexed{i: i, res: BulkCreateItem{Input: in, Doc: doc, Replay: replay, Err: err}}
		}()
	}
	for range inputs {
		r := <-out
		results[r.i] = r.res
	}
	return results
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
	saved, err := s.docs.Update(ctx, nil, result)
	if err != nil {
		return docEntity.Document{}, err
	}
	s.publishStatusEvent(ctx, saved)
	return saved, nil
}

func (s *service) publishStatusEvent(ctx context.Context, d docEntity.Document) {
	var err error
	switch d.Status {
	case enums.DocumentStatusGenerated:
		err = s.publisher.PublishDocumentGenerated(ctx, d)
	case enums.DocumentStatusFailed:
		err = s.publisher.PublishDocumentFailed(ctx, d)
	case enums.DocumentStatusCancelled:
		err = s.publisher.PublishDocumentCancelled(ctx, d)
	}
	if err != nil {
		log.Printf("documents: publish %s event: %v", d.Status, err)
	}
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
	saved, err := s.docs.Update(ctx, nil, result)
	if err != nil {
		return docEntity.Document{}, err
	}
	if pubErr := s.publisher.PublishDocumentCancelled(ctx, saved); pubErr != nil {
		log.Printf("documents: PublishDocumentCancelled: %v", pubErr)
	}
	return saved, nil
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
	if pubErr := s.publisher.PublishDocumentProcess(ctx, updated); pubErr != nil {
		log.Printf("documents: PublishDocumentProcess (retry): %v", pubErr)
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
	if s.storage != nil {
		url, err := s.storage.PresignedURL(ctx, *d.FilePath, 15*time.Minute)
		if err != nil {
			return "", err
		}
		return url, nil
	}
	return *d.FilePath, nil
}

// Preview merender template version dengan payload yang diberikan tanpa menyimpan ke DB.
func (s *service) Preview(ctx context.Context, templateID, versionID int64, tenantID *string, payload map[string]any) ([]byte, string, error) {
	tpl, err := s.templates.GetByID(ctx, nil, templateID, tenantID)
	if err != nil {
		return nil, "", mapRepoErr(err)
	}
	ver, err := s.versions.GetByID(ctx, nil, templateID, versionID, tenantID)
	if err != nil {
		return nil, "", mapRepoErr(err)
	}
	if len(ver.Schema) > 0 {
		if err := validators.ValidateSchema(ver.Schema, payload); err != nil {
			return nil, "", err
		}
	}
	gen := s.selector.Select(string(ver.OutputFormat), string(tpl.Engine))
	data, contentType, err := gen.Generate(ctx, ver.Content, payload)
	if err != nil {
		return nil, "", err
	}
	return data, contentType, nil
}

// Process dijalankan oleh Kafka consumer: QUEUED → PROCESSING → GENERATED (atau FAILED).
func (s *service) Process(ctx context.Context, id int64, tenantID *string) error {
	doc, err := s.docs.GetByID(ctx, nil, id, tenantID)
	if err != nil {
		return mapRepoErr(err)
	}
	if doc.Status != enums.DocumentStatusQueued {
		// Sudah diproses oleh consumer lain (at-least-once delivery) — skip.
		return nil
	}

	// Transisi QUEUED → PROCESSING
	processing, err := s.transitionDocument(ctx, doc, enums.DocumentStatusProcessing)
	if err != nil {
		return err
	}
	if _, err := s.docs.Update(ctx, nil, processing); err != nil {
		return err
	}

	// Transisi PROCESSING → GENERATED (toGenerated handler melakukan render file)
	generated, err := s.transitionDocument(ctx, processing, enums.DocumentStatusGenerated)
	if err != nil {
		// applyStateMachine sudah simpan status FAILED dan publish event Failed
		return err
	}
	saved, err := s.docs.Update(ctx, nil, generated)
	if err != nil {
		return err
	}
	if pubErr := s.publisher.PublishDocumentGenerated(ctx, saved); pubErr != nil {
		log.Printf("documents: Process: PublishDocumentGenerated: %v", pubErr)
	}
	return nil
}

// ZipDocuments mengunduh file tiap dokumen, membuat ZIP, dan mengembalikan URL download.
func (s *service) ZipDocuments(ctx context.Context, ids []int64, tenantID *string, label string) (string, error) {
	if s.storage == nil {
		return "", errors.New("storage provider not configured")
	}
	if len(ids) == 0 {
		return "", apperror.ErrInvalidInput
	}
	entries := make([]sharedStorage.ZipEntry, 0, len(ids))
	for _, id := range ids {
		d, err := s.docs.GetByID(ctx, nil, id, tenantID)
		if err != nil {
			return "", mapRepoErr(err)
		}
		if d.Status != enums.DocumentStatusGenerated || d.FilePath == nil {
			return "", fmt.Errorf("document %d belum generated", id)
		}
		data, err := s.storage.Download(ctx, *d.FilePath)
		if err != nil {
			return "", fmt.Errorf("download document %d: %w", id, err)
		}
		name := fmt.Sprintf("%d", d.ID)
		if d.FileName != nil {
			name = *d.FileName
		}
		entries = append(entries, sharedStorage.ZipEntry{Name: name, Data: data})
	}
	reqID := label
	if reqID == "" {
		reqID = fmt.Sprintf("zip-%d-docs", len(ids))
	}
	path, _, err := s.storage.Zip(ctx, 0, reqID, entries)
	if err != nil {
		return "", err
	}
	url, err := s.storage.PresignedURL(ctx, path, 15*time.Minute)
	if err != nil {
		return "", err
	}
	if pubErr := s.publisher.PublishDocumentsZipped(ctx, ids, tenantID, path, "zip"); pubErr != nil {
		log.Printf("documents: PublishDocumentsZipped: %v", pubErr)
	}
	return url, nil
}

// MergeDocuments menggabungkan file dokumen berformat sama menjadi satu, mengembalikan URL download.
// Untuk format teks (HTML, CSV): byte concat. Untuk PDF: butuh library pdfcpu — saat ini byte concat.
func (s *service) MergeDocuments(ctx context.Context, ids []int64, tenantID *string, label string) (string, error) {
	if s.storage == nil {
		return "", errors.New("storage provider not configured")
	}
	if len(ids) < 2 {
		return "", apperror.ErrInvalidInput
	}
	var format enums.OutputFormat
	srcPaths := make([]string, 0, len(ids))
	for i, id := range ids {
		d, err := s.docs.GetByID(ctx, nil, id, tenantID)
		if err != nil {
			return "", mapRepoErr(err)
		}
		if d.Status != enums.DocumentStatusGenerated || d.FilePath == nil {
			return "", fmt.Errorf("document %d belum generated", id)
		}
		if i == 0 {
			format = d.OutputFormat
		} else if d.OutputFormat != format {
			return "", fmt.Errorf("document %d format %s tidak cocok dengan %s", id, d.OutputFormat, format)
		}
		srcPaths = append(srcPaths, *d.FilePath)
	}
	ext := sharedStorage.ExtensionForFormat(string(format))
	reqID := label
	if reqID == "" {
		reqID = fmt.Sprintf("merge-%d-docs", len(ids))
	}
	path, _, err := s.storage.Compose(ctx, 0, reqID, srcPaths, ext)
	if err != nil {
		return "", err
	}
	url, err := s.storage.PresignedURL(ctx, path, 15*time.Minute)
	if err != nil {
		return "", err
	}
	if pubErr := s.publisher.PublishDocumentsMerged(ctx, ids, tenantID, path, string(format)); pubErr != nil {
		log.Printf("documents: PublishDocumentsMerged: %v", pubErr)
	}
	return url, nil
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
