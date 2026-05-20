package model

import (
	"time"

	cbEntity "go-document-generator/internal/entity/documentcallbackattempts"
)

type DocumentCallbackAttempt struct {
	ID                 int64          `gorm:"primaryKey;column:id"`
	DocumentID         int64          `gorm:"column:document_id"`
	CallbackURL        string         `gorm:"column:callback_url"`
	RequestPayload     map[string]any `gorm:"column:request_payload;serializer:json;type:jsonb"`
	ResponsePayload    map[string]any `gorm:"column:response_payload;serializer:json;type:jsonb"`
	ResponseStatusCode *int           `gorm:"column:response_status_code"`
	IsSuccess          bool           `gorm:"column:is_success"`
	ErrorMessage       *string        `gorm:"column:error_message"`
	AttemptedAt        time.Time      `gorm:"column:attempted_at"`
}

func (DocumentCallbackAttempt) TableName() string { return "document_callback_attempts" }

func ToEntity(m *DocumentCallbackAttempt) cbEntity.CallbackAttempt {
	if m == nil {
		return cbEntity.CallbackAttempt{}
	}
	return cbEntity.CallbackAttempt{
		ID:                 m.ID,
		DocumentID:         m.DocumentID,
		CallbackURL:        m.CallbackURL,
		RequestPayload:     m.RequestPayload,
		ResponsePayload:    m.ResponsePayload,
		ResponseStatusCode: m.ResponseStatusCode,
		IsSuccess:          m.IsSuccess,
		ErrorMessage:       m.ErrorMessage,
		AttemptedAt:        m.AttemptedAt,
	}
}

func ToModel(e cbEntity.CallbackAttempt) DocumentCallbackAttempt {
	return DocumentCallbackAttempt{
		ID:                 e.ID,
		DocumentID:         e.DocumentID,
		CallbackURL:        e.CallbackURL,
		RequestPayload:     e.RequestPayload,
		ResponsePayload:    e.ResponsePayload,
		ResponseStatusCode: e.ResponseStatusCode,
		IsSuccess:          e.IsSuccess,
		ErrorMessage:       e.ErrorMessage,
		AttemptedAt:        e.AttemptedAt,
	}
}
