package enums

type TemplateEngine string

const (
	TemplateEngineHandlebars TemplateEngine = "HANDLEBARS"
	TemplateEngineMustache   TemplateEngine = "MUSTACHE"
	TemplateEngineHTML       TemplateEngine = "HTML"
)

type OutputFormat string

const (
	OutputFormatPDF  OutputFormat = "PDF"
	OutputFormatHTML OutputFormat = "HTML"
	OutputFormatDOCX OutputFormat = "DOCX"
)

type DocumentStatus string

const (
	DocumentStatusPending    DocumentStatus = "PENDING"
	DocumentStatusQueued     DocumentStatus = "QUEUED"
	DocumentStatusProcessing DocumentStatus = "PROCESSING"
	DocumentStatusGenerated  DocumentStatus = "GENERATED"
	DocumentStatusFailed     DocumentStatus = "FAILED"
	DocumentStatusCancelled  DocumentStatus = "CANCELLED"
)

type DmsStatus string

const (
	DmsStatusNotSent DmsStatus = "NOT_SENT"
	DmsStatusQueued  DmsStatus = "QUEUED"
	DmsStatusSent    DmsStatus = "SENT"
	DmsStatusFailed  DmsStatus = "FAILED"
)

type CallbackStatus string

const (
	CallbackStatusPending  CallbackStatus = "PENDING"
	CallbackStatusSuccess  CallbackStatus = "SUCCESS"
	CallbackStatusFailed   CallbackStatus = "FAILED"
	CallbackStatusRetrying CallbackStatus = "RETRYING"
)

type StorageProvider string

const (
	StorageProviderLocal StorageProvider = "LOCAL"
	StorageProviderS3    StorageProvider = "S3"
	StorageProviderMinio StorageProvider = "MINIO"
	StorageProviderGCS   StorageProvider = "GCS"
	StorageProviderAzure StorageProvider = "AZURE"
)
