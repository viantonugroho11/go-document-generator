package model

import (
	"time"
)

type DocumentTemplateVersion struct {
	ID            int64             `gorm:"column:id;primaryKey"`
	TemplateID    int64             `gorm:"column:template_id"`
	Version       int               `gorm:"column:version"`
	Content       string            `gorm:"column:content"`
	Schema        map[string]any    `gorm:"column:schema;serializer:json"`
	SamplePayload map[string]any    `gorm:"column:sample_payload;serializer:json"`
	IsPublished   bool              `gorm:"column:is_published"`
	PublishedAt   *time.Time        `gorm:"column:published_at"`
	CreatedAt     time.Time         `gorm:"column:created_at"`
}

func (DocumentTemplateVersion) TableName() string {
	return "document_template_versions"
}

