package model

import (
	"time"

	verEntity "go-document-generator/internal/entity/documenttemplateversions"
	"go-document-generator/internal/entity/enums"
)

type DocumentTemplateVersion struct {
	ID            int64              `gorm:"primaryKey;column:id"`
	TenantID      *string            `gorm:"column:tenant_id;type:uuid"`
	TemplateID    int64              `gorm:"column:template_id"`
	Version       int                `gorm:"column:version"`
	Content       string             `gorm:"column:content"`
	Schema        map[string]any     `gorm:"column:schema;serializer:json;type:jsonb"`
	Variables     []any              `gorm:"column:variables;serializer:json;type:jsonb"`
	SamplePayload map[string]any     `gorm:"column:sample_payload;serializer:json;type:jsonb"`
	OutputFormat  enums.OutputFormat `gorm:"column:output_format;type:output_format"`
	Checksum      *string            `gorm:"column:checksum"`
	IsPublished   bool               `gorm:"column:is_published"`
	PublishedAt   *time.Time         `gorm:"column:published_at"`
	CreatedBy     *string            `gorm:"column:created_by"`
	CreatedAt     time.Time          `gorm:"column:created_at"`
}

func (DocumentTemplateVersion) TableName() string { return "document_template_versions" }

func ToEntity(m *DocumentTemplateVersion) verEntity.TemplateVersion {
	if m == nil {
		return verEntity.TemplateVersion{}
	}
	return verEntity.TemplateVersion{
		ID:            m.ID,
		TenantID:      m.TenantID,
		TemplateID:    m.TemplateID,
		Version:       m.Version,
		Content:       m.Content,
		Schema:        m.Schema,
		Variables:     m.Variables,
		SamplePayload: m.SamplePayload,
		OutputFormat:  m.OutputFormat,
		Checksum:      m.Checksum,
		IsPublished:   m.IsPublished,
		PublishedAt:   m.PublishedAt,
		CreatedBy:     m.CreatedBy,
		CreatedAt:     m.CreatedAt,
	}
}

func ToModel(e verEntity.TemplateVersion) DocumentTemplateVersion {
	return DocumentTemplateVersion{
		ID:            e.ID,
		TenantID:      e.TenantID,
		TemplateID:    e.TemplateID,
		Version:       e.Version,
		Content:       e.Content,
		Schema:        e.Schema,
		Variables:     e.Variables,
		SamplePayload: e.SamplePayload,
		OutputFormat:  e.OutputFormat,
		Checksum:      e.Checksum,
		IsPublished:   e.IsPublished,
		PublishedAt:   e.PublishedAt,
		CreatedBy:     e.CreatedBy,
		CreatedAt:     e.CreatedAt,
	}
}
