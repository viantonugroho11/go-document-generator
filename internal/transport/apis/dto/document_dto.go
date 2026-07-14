package dto

import (
	"time"

	docEntity "go-document-generator/internal/entity/documents"
	cbEntity "go-document-generator/internal/entity/documentcallbackattempts"
	logEntity "go-document-generator/internal/entity/documentrenderlogs"
	"go-document-generator/internal/entity/enums"
	ucDoc "go-document-generator/internal/usecase/documents"
)

type GeneratedDocumentResponse struct {
	ID                int64                  `json:"id"`
	TenantID          *string                `json:"tenant_id"`
	RequestID         string                 `json:"request_id"`
	TemplateID        *int64                 `json:"template_id"`
	TemplateVersionID *int64                 `json:"template_version_id"`
	TemplateCode      string                 `json:"template_code"`
	TemplateVersion   int                    `json:"template_version"`
	Payload           map[string]any         `json:"payload"`
	Metadata          map[string]any         `json:"metadata"`
	Status            enums.DocumentStatus   `json:"status"`
	ErrorMessage      *string                `json:"error_message"`
	OutputFormat      enums.OutputFormat     `json:"output_format"`
	FileName          *string                `json:"file_name"`
	FilePath          *string                `json:"file_path"`
	StorageProvider   *enums.StorageProvider `json:"storage_provider"`
	FileSize          *int64                 `json:"file_size"`
	Checksum          *string                `json:"checksum"`
	ContentType       *string                `json:"content_type"`
	IsSigned          bool                   `json:"is_signed"`
	SignatureProvider *string                `json:"signature_provider"`
	SignedAt          *time.Time             `json:"signed_at"`
	StoreToDms        bool                   `json:"store_to_dms"`
	DmsDocumentID     *string                `json:"dms_document_id"`
	DmsStatus         enums.DmsStatus        `json:"dms_status"`
	HasCallback       bool                   `json:"has_callback"`
	CallbackURL       *string                `json:"callback_url"`
	CallbackStatus    enums.CallbackStatus   `json:"callback_status"`
	CallbackLastAt    *time.Time             `json:"callback_last_at"`
	RetryCount        int                    `json:"retry_count"`
	NextRetryAt       *time.Time             `json:"next_retry_at"`
	ExpiredAt         *time.Time             `json:"expired_at"`
	CreatedBy         *string                `json:"created_by"`
	CreatedAt         time.Time              `json:"created_at"`
	ProcessedAt       *time.Time             `json:"processed_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	DeletedAt         *time.Time             `json:"deleted_at,omitempty"`
}

type CreateDocumentRequest struct {
	TenantID        *string            `json:"tenant_id"`
	RequestID       string             `json:"request_id"`
	TemplateCode    string             `json:"template_code"`
	TemplateVersion *int               `json:"template_version"`
	OutputFormat    enums.OutputFormat `json:"output_format"`
	Payload         map[string]any     `json:"payload"`
	Metadata        map[string]any     `json:"metadata"`
	StoreToDms      bool               `json:"store_to_dms"`
	HasCallback     bool               `json:"has_callback"`
	CallbackURL     *string            `json:"callback_url"`
	ExpiredAt       *time.Time         `json:"expired_at"`
	CreatedBy       *string            `json:"created_by"`
}

type PatchDocumentRequest struct {
	Status       *enums.DocumentStatus `json:"status"`
	Payload      map[string]any        `json:"payload"`
	Metadata     map[string]any        `json:"metadata"`
	OutputFormat *enums.OutputFormat   `json:"output_format"`
	StoreToDms   *bool                 `json:"store_to_dms"`
	HasCallback  *bool                 `json:"has_callback"`
	CallbackURL  *string               `json:"callback_url"`
	ExpiredAt    *time.Time            `json:"expired_at"`
	ErrorMessage *string               `json:"error_message"`
}

type DocumentListResponse struct {
	Data []GeneratedDocumentResponse `json:"data"`
	Meta PaginationMeta              `json:"meta"`
}

type RenderLogResponse struct {
	ID              int64                `json:"id"`
	DocumentID      int64                `json:"document_id"`
	Status          enums.DocumentStatus `json:"status"`
	Message         *string              `json:"message"`
	ExecutionTimeMs *int64               `json:"execution_time_ms"`
	StackTrace      *string              `json:"stack_trace"`
	WorkerName      *string              `json:"worker_name"`
	CreatedAt       time.Time            `json:"created_at"`
}

type RenderLogListResponse struct {
	Data []RenderLogResponse `json:"data"`
	Meta PaginationMeta      `json:"meta"`
}

type CallbackAttemptResponse struct {
	ID                 int64          `json:"id"`
	DocumentID         int64          `json:"document_id"`
	CallbackURL        string         `json:"callback_url"`
	RequestPayload     map[string]any `json:"request_payload"`
	ResponsePayload    map[string]any `json:"response_payload"`
	ResponseStatusCode *int           `json:"response_status_code"`
	IsSuccess          bool           `json:"is_success"`
	ErrorMessage       *string        `json:"error_message"`
	AttemptedAt        time.Time      `json:"attempted_at"`
}

type CallbackAttemptListResponse struct {
	Data []CallbackAttemptResponse `json:"data"`
	Meta PaginationMeta            `json:"meta"`
}

type TestCallbackRequest struct {
	CallbackURL   string         `json:"callback_url"`
	SamplePayload map[string]any `json:"sample_payload"`
}

type TestCallbackResponse struct {
	Success            bool   `json:"success"`
	ResponseStatusCode int    `json:"response_status_code,omitempty"`
	ErrorMessage       string `json:"error_message,omitempty"`
}

func DocumentFromEntity(d docEntity.Document) GeneratedDocumentResponse {
	return GeneratedDocumentResponse{
		ID: d.ID, TenantID: d.TenantID, RequestID: d.RequestID,
		TemplateID: d.TemplateID, TemplateVersionID: d.TemplateVersionID,
		TemplateCode: d.TemplateCode, TemplateVersion: d.TemplateVersion,
		Payload: d.Payload, Metadata: d.Metadata, Status: d.Status, ErrorMessage: d.ErrorMessage,
		OutputFormat: d.OutputFormat, FileName: d.FileName, FilePath: d.FilePath,
		StorageProvider: d.StorageProvider, FileSize: d.FileSize, Checksum: d.Checksum,
		ContentType: d.ContentType, IsSigned: d.IsSigned, SignatureProvider: d.SignatureProvider,
		SignedAt: d.SignedAt, StoreToDms: d.StoreToDms, DmsDocumentID: d.DmsDocumentID,
		DmsStatus: d.DmsStatus, HasCallback: d.HasCallback, CallbackURL: d.CallbackURL,
		CallbackStatus: d.CallbackStatus, CallbackLastAt: d.CallbackLastAt,
		RetryCount: d.RetryCount, NextRetryAt: d.NextRetryAt, ExpiredAt: d.ExpiredAt,
		CreatedBy: d.CreatedBy, CreatedAt: d.CreatedAt, ProcessedAt: d.ProcessedAt,
		UpdatedAt: d.UpdatedAt, DeletedAt: d.DeletedAt,
	}
}

func ApplyPatchDocument(existing docEntity.Document, req PatchDocumentRequest) docEntity.Document {
	if req.Status != nil {
		existing.Status = *req.Status
	}
	if req.ErrorMessage != nil {
		existing.ErrorMessage = req.ErrorMessage
	}
	if req.Payload != nil {
		existing.Payload = req.Payload
	}
	if req.Metadata != nil {
		existing.Metadata = req.Metadata
	}
	if req.OutputFormat != nil {
		existing.OutputFormat = *req.OutputFormat
	}
	if req.StoreToDms != nil {
		existing.StoreToDms = *req.StoreToDms
	}
	if req.HasCallback != nil {
		existing.HasCallback = *req.HasCallback
	}
	if req.CallbackURL != nil {
		existing.CallbackURL = req.CallbackURL
	}
	if req.ExpiredAt != nil {
		existing.ExpiredAt = req.ExpiredAt
	}
	return existing
}

func (r CreateDocumentRequest) ToInput(headerTenant *string) ucDoc.CreateInput {
	tid := ResolveTenant(headerTenant, r.TenantID)
	return ucDoc.CreateInput{
		TenantID: tid, RequestID: r.RequestID, TemplateCode: r.TemplateCode,
		TemplateVersion: r.TemplateVersion, OutputFormat: r.OutputFormat,
		Payload: r.Payload, Metadata: r.Metadata, StoreToDms: r.StoreToDms,
		HasCallback: r.HasCallback, CallbackURL: r.CallbackURL,
		ExpiredAt: r.ExpiredAt, CreatedBy: r.CreatedBy,
	}
}

func RenderLogFromEntity(l logEntity.RenderLog) RenderLogResponse {
	return RenderLogResponse{
		ID: l.ID, DocumentID: l.DocumentID, Status: l.Status, Message: l.Message,
		ExecutionTimeMs: l.ExecutionTimeMs, StackTrace: l.StackTrace,
		WorkerName: l.WorkerName, CreatedAt: l.CreatedAt,
	}
}

func CallbackFromEntity(a cbEntity.CallbackAttempt) CallbackAttemptResponse {
	return CallbackAttemptResponse{
		ID: a.ID, DocumentID: a.DocumentID, CallbackURL: a.CallbackURL,
		RequestPayload: a.RequestPayload, ResponsePayload: a.ResponsePayload,
		ResponseStatusCode: a.ResponseStatusCode, IsSuccess: a.IsSuccess,
		ErrorMessage: a.ErrorMessage, AttemptedAt: a.AttemptedAt,
	}
}

// --- Bulk Create ---

type BulkCreateDocumentRequest struct {
	Items []CreateDocumentRequest `json:"items"`
}

type BulkCreateDocumentItemResponse struct {
	RequestID string                     `json:"request_id"`
	Replay    bool                       `json:"replay,omitempty"`
	Doc       *GeneratedDocumentResponse `json:"doc,omitempty"`
	Error     string                     `json:"error,omitempty"`
}

type BulkCreateDocumentResponse struct {
	Items     []BulkCreateDocumentItemResponse `json:"items"`
	Total     int                              `json:"total"`
	Succeeded int                              `json:"succeeded"`
	Failed    int                              `json:"failed"`
}

// --- Preview ---

type PreviewDocumentRequest struct {
	Payload map[string]any `json:"payload"`
}
