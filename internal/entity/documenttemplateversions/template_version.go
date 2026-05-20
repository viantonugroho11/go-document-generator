package documenttemplateversions

import (
	"time"

	"go-document-generator/internal/entity/enums"
)

type TemplateVersion struct {
	ID            int64
	TenantID      *string
	TemplateID    int64
	Version       int
	Content       string
	Schema        map[string]any
	Variables     []any
	SamplePayload map[string]any
	OutputFormat  enums.OutputFormat
	Checksum      *string
	IsPublished   bool
	PublishedAt   *time.Time
	CreatedBy     *string
	CreatedAt     time.Time
}
