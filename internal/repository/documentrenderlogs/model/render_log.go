package model

import (
	"time"

	logEntity "go-document-generator/internal/entity/documentrenderlogs"
	"go-document-generator/internal/entity/enums"
)

type DocumentRenderLog struct {
	ID              int64                `gorm:"primaryKey;column:id"`
	DocumentID      int64                `gorm:"column:document_id"`
	Status          enums.DocumentStatus `gorm:"column:status;type:document_status"`
	Message         *string              `gorm:"column:message"`
	ExecutionTimeMs *int64               `gorm:"column:execution_time_ms"`
	StackTrace      *string              `gorm:"column:stack_trace"`
	WorkerName      *string              `gorm:"column:worker_name"`
	CreatedAt       time.Time            `gorm:"column:created_at"`
}

func (DocumentRenderLog) TableName() string { return "document_render_logs" }

func ToEntity(m *DocumentRenderLog) logEntity.RenderLog {
	if m == nil {
		return logEntity.RenderLog{}
	}
	return logEntity.RenderLog{
		ID:              m.ID,
		DocumentID:      m.DocumentID,
		Status:          m.Status,
		Message:         m.Message,
		ExecutionTimeMs: m.ExecutionTimeMs,
		StackTrace:      m.StackTrace,
		WorkerName:      m.WorkerName,
		CreatedAt:       m.CreatedAt,
	}
}

func ToModel(e logEntity.RenderLog) DocumentRenderLog {
	return DocumentRenderLog{
		ID:              e.ID,
		DocumentID:      e.DocumentID,
		Status:          e.Status,
		Message:         e.Message,
		ExecutionTimeMs: e.ExecutionTimeMs,
		StackTrace:      e.StackTrace,
		WorkerName:      e.WorkerName,
		CreatedAt:       e.CreatedAt,
	}
}
