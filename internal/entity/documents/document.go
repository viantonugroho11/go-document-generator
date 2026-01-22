package documents

import "time"

type Document struct {
	ID              int64          `json:"id"`
	RequestID       string         `json:"request_id"`
	TemplateCode    string         `json:"template_code"`
	TemplateVersion *int           `json:"template_version"`
	Payload         map[string]any `json:"payload"`
	Metadata        map[string]any `json:"metadata"`
	Status          string         `json:"status"`
	ErrorMessage    *string        `json:"error_message"`
	FileName        *string        `json:"file_name"`
	FilePath        *string        `json:"file_path"`
	FileSize        *int64         `json:"file_size"`
	Checksum        *string        `json:"checksum"`
	ContentType     *string        `json:"content_type"`
	StoreToDMS      bool           `json:"store_to_dms"`
	DMSDocumentID   *string        `json:"dms_document_id"`
	DMSStatus       string         `json:"dms_status"`
	HasCallback     bool           `json:"has_callback"`
	CallbackURL     *string        `json:"callback_url"`
	CallbackStatus  string         `json:"callback_status"`
	CreatedBy       *string        `json:"created_by"`
	CreatedAt       time.Time      `json:"created_at"`
	ProcessedAt     *time.Time     `json:"processed_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}
