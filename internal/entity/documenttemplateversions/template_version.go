package documenttemplateversions

import "time"

type DocumentTemplateVersion struct {
	ID            int64                  `json:"id"`
	TemplateID    int64                  `json:"template_id"`
	Version       int                    `json:"version"`
	Content       string                 `json:"content"`
	Schema        map[string]any         `json:"schema"`
	SamplePayload map[string]any         `json:"sample_payload"`
	IsPublished   bool                   `json:"is_published"`
	PublishedAt   *time.Time             `json:"published_at"`
	CreatedAt     time.Time              `json:"created_at"`
}

