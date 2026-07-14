package events

import "time"

const MessageSchemaVersion = 1

// EventMeta metadata standar untuk semua event.
type EventMeta struct {
	EventID              string `json:"event_id"`
	EventTimestamp       string `json:"event_timestamp"`        // RFC3339 UTC
	Action               string `json:"action"`                 // CREATE | UPDATE
	Resource             string `json:"resource"`               // "Document" | "DocumentBulk"
	MessageSchemaVersion int    `json:"message_schema_version"` // selalu 1 untuk sekarang
}

// DocumentState snapshot dokumen pada satu titik waktu.
type DocumentState struct {
	ID              int64      `json:"id"`
	RequestID       string     `json:"request_id"`
	TenantID        *string    `json:"tenant_id,omitempty"`
	TemplateCode    string     `json:"template_code"`
	TemplateVersion int        `json:"template_version"`
	OutputFormat    string     `json:"output_format"`
	Status          string     `json:"status"`
	ErrorMessage    *string    `json:"error_message,omitempty"`
	RetryCount      int        `json:"retry_count"`
	FilePath        *string    `json:"file_path,omitempty"`
	FileName        *string    `json:"file_name,omitempty"`
	FileSize        *int64     `json:"file_size,omitempty"`
	ContentType     *string    `json:"content_type,omitempty"`
	Checksum        *string    `json:"checksum,omitempty"`
	StorageProvider *string    `json:"storage_provider,omitempty"`
	ProcessedAt     *time.Time `json:"processed_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// DocumentEvent envelope untuk semua lifecycle event dokumen.
// resource_id = request_id (natural business key, dipakai sebagai partition key).
// before = nil saat action CREATE.
type DocumentEvent struct {
	ResourceID string         `json:"resource_id"`
	Meta       EventMeta      `json:"meta"`
	Before     *DocumentState `json:"before"`
	After      *DocumentState `json:"after"`
}

// DocumentBulkState snapshot operasi bulk (zip / merge).
type DocumentBulkState struct {
	DocumentIDs  []int64 `json:"document_ids"`
	TenantID     *string `json:"tenant_id,omitempty"`
	OutputPath   string  `json:"output_path"`
	OutputFormat string  `json:"output_format,omitempty"`
}

// DocumentBulkEvent envelope untuk operasi zip dan merge.
// resource_id = output_path.
// before selalu nil (bulk operation selalu CREATE).
type DocumentBulkEvent struct {
	ResourceID string             `json:"resource_id"`
	Meta       EventMeta          `json:"meta"`
	Before     *DocumentBulkState `json:"before"`
	After      *DocumentBulkState `json:"after"`
}
