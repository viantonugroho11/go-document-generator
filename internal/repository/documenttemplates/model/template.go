package model

import (
	entity "go-document-generator/internal/entity/documenttemplates"
	"time"
)

type DocumentTemplate struct {
	ID           int64      `gorm:"column:id;primaryKey"`
	Code         string     `gorm:"column:code"`
	Name         string     `gorm:"column:name"`
	Description  *string    `gorm:"column:description"`
	Engine       string     `gorm:"column:engine"`
	OutputFormat string     `gorm:"column:output_format"`
	IsActive     bool       `gorm:"column:is_active"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
}

func (DocumentTemplate) TableName() string {
	return "document_templates"
}

func (m *DocumentTemplate) ToEntity() entity.DocumentTemplate {
	return entity.DocumentTemplate{
		ID:           m.ID,
		Code:         m.Code,
		Name:         m.Name,
		Description:  m.Description,
		Engine:       m.Engine,
		OutputFormat: m.OutputFormat,
		IsActive:     m.IsActive,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func FromEntity(e entity.DocumentTemplate) DocumentTemplate {
	return DocumentTemplate{
		ID:           e.ID,
		Code:         e.Code,
		Name:         e.Name,
		Description:  e.Description,
		Engine:       e.Engine,
		OutputFormat: e.OutputFormat,
		IsActive:     e.IsActive,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}