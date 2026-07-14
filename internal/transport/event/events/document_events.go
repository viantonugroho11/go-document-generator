package events

// DocumentQueuedEvent payload Kafka saat dokumen masuk antrian render.
type DocumentQueuedEvent struct {
	ID           int64  `json:"id"`
	RequestID    string `json:"request_id"`
	TemplateCode string `json:"template_code"`
	Status       string `json:"status"`
}

// DocumentRetriedEvent payload Kafka saat retry generation.
type DocumentRetriedEvent struct {
	ID        int64  `json:"id"`
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
}

// DocumentGeneratedEvent payload Kafka saat dokumen berhasil di-render.
type DocumentGeneratedEvent struct {
	ID          int64   `json:"id"`
	RequestID   string  `json:"request_id"`
	Status      string  `json:"status"`
	OutputFormat string `json:"output_format"`
	FilePath    *string `json:"file_path,omitempty"`
	FileSize    *int64  `json:"file_size,omitempty"`
}

// DocumentFailedEvent payload Kafka saat render gagal.
type DocumentFailedEvent struct {
	ID           int64   `json:"id"`
	RequestID    string  `json:"request_id"`
	Status       string  `json:"status"`
	ErrorMessage *string `json:"error_message,omitempty"`
}

// DocumentCancelledEvent payload Kafka saat dokumen dibatalkan.
type DocumentCancelledEvent struct {
	ID        int64  `json:"id"`
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
}
