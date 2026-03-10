package model

import (
	entity "go-document-generator/internal/entity/documenttemplateversions"
	"time"
)

type DocumentTemplateVersion struct {
	ID         int64  `gorm:"column:id;primaryKey"`
	TemplateID int64  `gorm:"column:template_id"`
	Version    int    `gorm:"column:version"`
	Content    string `gorm:"column:content"`
	// file path
	// FileTemplatePath string         `gorm:"column:file_template_path"` // file template path
	Schema           map[string]any `gorm:"column:schema;serializer:json"`
	SamplePayload    map[string]any `gorm:"column:sample_payload;serializer:json"`
	IsPublished      bool           `gorm:"column:is_published"`
	PublishedAt      *time.Time     `gorm:"column:published_at"`
	CreatedAt        time.Time      `gorm:"column:created_at"`
}

func (DocumentTemplateVersion) TableName() string {
	return "document_template_versions"
}

func (m *DocumentTemplateVersion) ToEntity() entity.DocumentTemplateVersion {
	return entity.DocumentTemplateVersion{
		ID:            m.ID,
		TemplateID:    m.TemplateID,
		Version:       m.Version,
		Content:       m.Content,
		// FileTemplatePath: m.FileTemplatePath,
		Schema:        m.Schema,
		SamplePayload: m.SamplePayload,
		IsPublished:   m.IsPublished,
		PublishedAt:   m.PublishedAt,
	}
}

// to model from entity
func FromEntity(e entity.DocumentTemplateVersion) DocumentTemplateVersion {
	return DocumentTemplateVersion{
		ID:            e.ID,
		TemplateID:    e.TemplateID,
		Version:       e.Version,
		Content:       e.Content,
		Schema:        e.Schema,
		SamplePayload: e.SamplePayload,
		IsPublished:   e.IsPublished,
		PublishedAt:   e.PublishedAt,
		CreatedAt:     e.CreatedAt,
	}
}

// params for list
type DocumentTemplateVersionListParams struct {
	TemplateID  int64
	Version     int
	IsPublished bool
	Limit       int
	Offset      int
	OrderBy     string
	Ids         []string
}
