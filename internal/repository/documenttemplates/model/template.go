package model

import (
	"time"

	tplEntity "go-document-generator/internal/entity/documenttemplates"
	"go-document-generator/internal/entity/enums"
)

type DocumentTemplate struct {
	ID            int64              `gorm:"primaryKey;column:id"`
	TenantID      *string            `gorm:"column:tenant_id;type:uuid"`
	Code          string             `gorm:"column:code"`
	Name          string             `gorm:"column:name"`
	Description   *string            `gorm:"column:description"`
	Engine        enums.TemplateEngine `gorm:"column:engine;type:template_engine"`
	DefaultFormat enums.OutputFormat `gorm:"column:default_format;type:output_format"`
	Category      *string            `gorm:"column:category"`
	IsActive      bool               `gorm:"column:is_active"`
	CreatedBy     *string            `gorm:"column:created_by"`
	UpdatedBy     *string            `gorm:"column:updated_by"`
	CreatedAt     time.Time          `gorm:"column:created_at"`
	UpdatedAt     time.Time          `gorm:"column:updated_at"`
}

func (DocumentTemplate) TableName() string { return "document_templates" }

func ToEntity(m *DocumentTemplate) tplEntity.Template {
	if m == nil {
		return tplEntity.Template{}
	}
	return tplEntity.Template{
		ID:            m.ID,
		TenantID:      m.TenantID,
		Code:          m.Code,
		Name:          m.Name,
		Description:   m.Description,
		Engine:        m.Engine,
		DefaultFormat: m.DefaultFormat,
		Category:      m.Category,
		IsActive:      m.IsActive,
		CreatedBy:     m.CreatedBy,
		UpdatedBy:     m.UpdatedBy,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func ToModel(e tplEntity.Template) DocumentTemplate {
	return DocumentTemplate{
		ID:            e.ID,
		TenantID:      e.TenantID,
		Code:          e.Code,
		Name:          e.Name,
		Description:   e.Description,
		Engine:        e.Engine,
		DefaultFormat: e.DefaultFormat,
		Category:      e.Category,
		IsActive:      e.IsActive,
		CreatedBy:     e.CreatedBy,
		UpdatedBy:     e.UpdatedBy,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}
