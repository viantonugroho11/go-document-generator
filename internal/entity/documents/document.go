package documents

import (
	"time"

	"go-document-generator/internal/entity/enums"
)

type Document struct {
	ID                 int64
	TenantID           *string
	RequestID          string
	TemplateID         *int64
	TemplateVersionID  *int64
	TemplateCode       string
	TemplateVersion    int
	Payload            map[string]any
	Metadata           map[string]any
	Status             enums.DocumentStatus
	ErrorMessage       *string
	OutputFormat       enums.OutputFormat
	FileName           *string
	FilePath           *string
	StorageProvider    *enums.StorageProvider
	FileSize           *int64
	Checksum           *string
	ContentType        *string
	IsSigned           bool
	SignatureProvider  *string
	SignedAt           *time.Time
	StoreToDms         bool
	DmsDocumentID      *string
	DmsStatus          enums.DmsStatus
	HasCallback        bool
	CallbackURL        *string
	CallbackStatus     enums.CallbackStatus
	CallbackLastAt     *time.Time
	RetryCount         int
	NextRetryAt        *time.Time
	ExpiredAt          *time.Time
	CreatedBy          *string
	CreatedAt          time.Time
	ProcessedAt        *time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time
}
