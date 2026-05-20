package documenttemplates

import (
	"time"

	"go-document-generator/internal/entity/enums"
)

type Template struct {
	ID            int64
	TenantID      *string
	Code          string
	Name          string
	Description   *string
	Engine        enums.TemplateEngine
	DefaultFormat enums.OutputFormat
	Category      *string
	IsActive      bool
	CreatedBy     *string
	UpdatedBy     *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
