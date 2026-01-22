package model

import "time"

type Document struct {
	ID              int64              `gorm:"column:id;primaryKey"`
	RequestID       string             `gorm:"column:request_id"`
	TemplateCode    string             `gorm:"column:template_code"`
	TemplateVersion *int               `gorm:"column:template_version"`
	Payload         map[string]any     `gorm:"column:payload;serializer:json"`
	Metadata        map[string]any     `gorm:"column:metadata;serializer:json"`
	Status          string             `gorm:"column:status"`
	ErrorMessage    *string            `gorm:"column:error_message"`
	FileName        *string            `gorm:"column:file_name"`
	FilePath        *string            `gorm:"column:file_path"`
	FileSize        *int64             `gorm:"column:file_size"`
	Checksum        *string            `gorm:"column:checksum"`
	ContentType     *string            `gorm:"column:content_type"`
	StoreToDMS      bool               `gorm:"column:store_to_dms"`
	DMSDocumentID   *string            `gorm:"column:dms_document_id"`
	DMSStatus       string             `gorm:"column:dms_status"`
	HasCallback     bool               `gorm:"column:has_callback"`
	CallbackURL     *string            `gorm:"column:callback_url"`
	CallbackStatus  string             `gorm:"column:callback_status"`
	CreatedBy       *string            `gorm:"column:created_by"`
	CreatedAt       time.Time          `gorm:"column:created_at"`
	ProcessedAt     *time.Time         `gorm:"column:processed_at"`
	UpdatedAt       time.Time          `gorm:"column:updated_at"`
}

func (Document) TableName() string {
	return "documents"
}

