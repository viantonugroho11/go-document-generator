package events

import "time"

// DocumentQueuedEvent dipublish ke document-events (observability)
// DAN document-process (trigger generation consumer).
type DocumentQueuedEvent struct {
	ID              int64      `json:"id"`
	RequestID       string     `json:"request_id"`
	TenantID        *string    `json:"tenant_id,omitempty"`
	TemplateCode    string     `json:"template_code"`
	TemplateVersion int        `json:"template_version"`
	OutputFormat    string     `json:"output_format"`
	Status          string     `json:"status"`
	OccurredAt      time.Time  `json:"occurred_at"`
}

// DocumentRetriedEvent dipublish ke document-events saat retry diminta.
type DocumentRetriedEvent struct {
	ID         int64     `json:"id"`
	RequestID  string    `json:"request_id"`
	TenantID   *string   `json:"tenant_id,omitempty"`
	RetryCount int       `json:"retry_count"`
	Status     string    `json:"status"`
	OccurredAt time.Time `json:"occurred_at"`
}

// DocumentGeneratedEvent dipublish ke document-events saat render berhasil.
type DocumentGeneratedEvent struct {
	ID              int64      `json:"id"`
	RequestID       string     `json:"request_id"`
	TenantID        *string    `json:"tenant_id,omitempty"`
	TemplateCode    string     `json:"template_code"`
	OutputFormat    string     `json:"output_format"`
	Status          string     `json:"status"`
	FileName        *string    `json:"file_name,omitempty"`
	FilePath        *string    `json:"file_path,omitempty"`
	FileSize        *int64     `json:"file_size,omitempty"`
	ContentType     *string    `json:"content_type,omitempty"`
	Checksum        *string    `json:"checksum,omitempty"`
	StorageProvider *string    `json:"storage_provider,omitempty"`
	ProcessedAt     *time.Time `json:"processed_at,omitempty"`
	OccurredAt      time.Time  `json:"occurred_at"`
}

// DocumentFailedEvent dipublish ke document-events saat render gagal.
type DocumentFailedEvent struct {
	ID           int64     `json:"id"`
	RequestID    string    `json:"request_id"`
	TenantID     *string   `json:"tenant_id,omitempty"`
	Status       string    `json:"status"`
	ErrorMessage *string   `json:"error_message,omitempty"`
	RetryCount   int       `json:"retry_count"`
	OccurredAt   time.Time `json:"occurred_at"`
}

// DocumentCancelledEvent dipublish ke document-events saat dokumen dibatalkan.
type DocumentCancelledEvent struct {
	ID         int64     `json:"id"`
	RequestID  string    `json:"request_id"`
	TenantID   *string   `json:"tenant_id,omitempty"`
	Status     string    `json:"status"`
	OccurredAt time.Time `json:"occurred_at"`
}

// DocumentsZippedEvent dipublish ke document-events saat operasi zip selesai.
type DocumentsZippedEvent struct {
	DocumentIDs []int64   `json:"document_ids"`
	TenantID    *string   `json:"tenant_id,omitempty"`
	ZipPath     string    `json:"zip_path"`
	OccurredAt  time.Time `json:"occurred_at"`
}

// DocumentsMergedEvent dipublish ke document-events saat operasi merge selesai.
type DocumentsMergedEvent struct {
	DocumentIDs  []int64   `json:"document_ids"`
	TenantID     *string   `json:"tenant_id,omitempty"`
	MergedPath   string    `json:"merged_path"`
	OutputFormat string    `json:"output_format"`
	OccurredAt   time.Time `json:"occurred_at"`
}
