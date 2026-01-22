package model

import "time"

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

