package events

// DocumentQueuedEvent payload Kafka saat dokumen masuk antrian render.
type DocumentQueuedEvent struct {
	ID         int64  `json:"id"`
	RequestID  string `json:"request_id"`
	TemplateCode string `json:"template_code"`
	Status     string `json:"status"`
}

// DocumentRetriedEvent payload Kafka saat retry generation.
type DocumentRetriedEvent struct {
	ID        int64  `json:"id"`
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
}
