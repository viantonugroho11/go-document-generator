package events

// TemplateCreatedEvent payload Kafka untuk template baru.
type TemplateCreatedEvent struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
}

// TemplateUpdatedEvent payload Kafka untuk perubahan template.
type TemplateUpdatedEvent struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
}

// TemplateVersionCreatedEvent payload Kafka untuk versi template baru.
type TemplateVersionCreatedEvent struct {
	ID         int64 `json:"id"`
	TemplateID int64 `json:"template_id"`
	Version    int   `json:"version"`
}

// TemplateVersionPublishedEvent payload Kafka saat versi dipublish.
type TemplateVersionPublishedEvent struct {
	ID         int64 `json:"id"`
	TemplateID int64 `json:"template_id"`
	Version    int   `json:"version"`
}
