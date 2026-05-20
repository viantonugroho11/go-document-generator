package documentrenderlogs

import (
	"time"

	"go-document-generator/internal/entity/enums"
)

type RenderLog struct {
	ID              int64
	DocumentID      int64
	Status          enums.DocumentStatus
	Message         *string
	ExecutionTimeMs *int64
	StackTrace      *string
	WorkerName      *string
	CreatedAt       time.Time
}
