package model

import (
	"time"

	docEntity "go-document-generator/internal/entity/documents"
	"go-document-generator/internal/entity/enums"
)

type Document struct {
	ID                int64                 `gorm:"primaryKey;column:id"`
	TenantID          *string               `gorm:"column:tenant_id;type:uuid"`
	RequestID         string                `gorm:"column:request_id"`
	TemplateID        *int64                `gorm:"column:template_id"`
	TemplateVersionID *int64                `gorm:"column:template_version_id"`
	TemplateCode      string                `gorm:"column:template_code"`
	TemplateVersion   int                   `gorm:"column:template_version"`
	Payload           map[string]any        `gorm:"column:payload;serializer:json;type:jsonb"`
	Metadata          map[string]any        `gorm:"column:metadata;serializer:json;type:jsonb"`
	Status            enums.DocumentStatus  `gorm:"column:status;type:document_status"`
	ErrorMessage      *string               `gorm:"column:error_message"`
	OutputFormat      enums.OutputFormat    `gorm:"column:output_format;type:output_format"`
	FileName          *string               `gorm:"column:file_name"`
	FilePath          *string               `gorm:"column:file_path"`
	StorageProvider   *enums.StorageProvider `gorm:"column:storage_provider;type:storage_provider"`
	FileSize          *int64                `gorm:"column:file_size"`
	Checksum          *string               `gorm:"column:checksum"`
	ContentType       *string               `gorm:"column:content_type"`
	IsSigned          bool                  `gorm:"column:is_signed"`
	SignatureProvider *string               `gorm:"column:signature_provider"`
	SignedAt          *time.Time            `gorm:"column:signed_at"`
	StoreToDms        bool                  `gorm:"column:store_to_dms"`
	DmsDocumentID     *string               `gorm:"column:dms_document_id"`
	DmsStatus         enums.DmsStatus       `gorm:"column:dms_status;type:dms_status"`
	HasCallback       bool                  `gorm:"column:has_callback"`
	CallbackURL       *string               `gorm:"column:callback_url"`
	CallbackStatus    enums.CallbackStatus  `gorm:"column:callback_status;type:callback_status"`
	CallbackLastAt    *time.Time            `gorm:"column:callback_last_at"`
	RetryCount        int                   `gorm:"column:retry_count"`
	NextRetryAt       *time.Time            `gorm:"column:next_retry_at"`
	ExpiredAt         *time.Time            `gorm:"column:expired_at"`
	CreatedBy         *string               `gorm:"column:created_by"`
	CreatedAt         time.Time             `gorm:"column:created_at"`
	ProcessedAt       *time.Time            `gorm:"column:processed_at"`
	UpdatedAt         time.Time             `gorm:"column:updated_at"`
	DeletedAt         *time.Time            `gorm:"column:deleted_at"`
}

func (Document) TableName() string { return "documents" }

func ToEntity(m *Document) docEntity.Document {
	if m == nil {
		return docEntity.Document{}
	}
	return docEntity.Document{
		ID:                m.ID,
		TenantID:          m.TenantID,
		RequestID:         m.RequestID,
		TemplateID:        m.TemplateID,
		TemplateVersionID: m.TemplateVersionID,
		TemplateCode:      m.TemplateCode,
		TemplateVersion:   m.TemplateVersion,
		Payload:           m.Payload,
		Metadata:          m.Metadata,
		Status:            m.Status,
		ErrorMessage:      m.ErrorMessage,
		OutputFormat:      m.OutputFormat,
		FileName:          m.FileName,
		FilePath:          m.FilePath,
		StorageProvider:   m.StorageProvider,
		FileSize:          m.FileSize,
		Checksum:          m.Checksum,
		ContentType:       m.ContentType,
		IsSigned:          m.IsSigned,
		SignatureProvider: m.SignatureProvider,
		SignedAt:          m.SignedAt,
		StoreToDms:        m.StoreToDms,
		DmsDocumentID:     m.DmsDocumentID,
		DmsStatus:         m.DmsStatus,
		HasCallback:       m.HasCallback,
		CallbackURL:       m.CallbackURL,
		CallbackStatus:    m.CallbackStatus,
		CallbackLastAt:    m.CallbackLastAt,
		RetryCount:        m.RetryCount,
		NextRetryAt:       m.NextRetryAt,
		ExpiredAt:         m.ExpiredAt,
		CreatedBy:         m.CreatedBy,
		CreatedAt:         m.CreatedAt,
		ProcessedAt:       m.ProcessedAt,
		UpdatedAt:         m.UpdatedAt,
		DeletedAt:         m.DeletedAt,
	}
}

func ToModel(e docEntity.Document) Document {
	return Document{
		ID:                e.ID,
		TenantID:          e.TenantID,
		RequestID:         e.RequestID,
		TemplateID:        e.TemplateID,
		TemplateVersionID: e.TemplateVersionID,
		TemplateCode:      e.TemplateCode,
		TemplateVersion:   e.TemplateVersion,
		Payload:           e.Payload,
		Metadata:          e.Metadata,
		Status:            e.Status,
		ErrorMessage:      e.ErrorMessage,
		OutputFormat:      e.OutputFormat,
		FileName:          e.FileName,
		FilePath:          e.FilePath,
		StorageProvider:   e.StorageProvider,
		FileSize:          e.FileSize,
		Checksum:          e.Checksum,
		ContentType:       e.ContentType,
		IsSigned:          e.IsSigned,
		SignatureProvider: e.SignatureProvider,
		SignedAt:          e.SignedAt,
		StoreToDms:        e.StoreToDms,
		DmsDocumentID:     e.DmsDocumentID,
		DmsStatus:         e.DmsStatus,
		HasCallback:       e.HasCallback,
		CallbackURL:       e.CallbackURL,
		CallbackStatus:    e.CallbackStatus,
		CallbackLastAt:    e.CallbackLastAt,
		RetryCount:        e.RetryCount,
		NextRetryAt:       e.NextRetryAt,
		ExpiredAt:         e.ExpiredAt,
		CreatedBy:         e.CreatedBy,
		CreatedAt:         e.CreatedAt,
		ProcessedAt:       e.ProcessedAt,
		UpdatedAt:         e.UpdatedAt,
		DeletedAt:         e.DeletedAt,
	}
}
