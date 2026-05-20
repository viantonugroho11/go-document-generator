package dto

import (
	"time"

	verEntity "go-document-generator/internal/entity/documenttemplateversions"
	"go-document-generator/internal/entity/enums"
)

type TemplateVersionResponse struct {
	ID            int64              `json:"id"`
	TenantID      *string            `json:"tenant_id"`
	TemplateID    int64              `json:"template_id"`
	Version       int                `json:"version"`
	Content       string             `json:"content,omitempty"`
	Schema        map[string]any     `json:"schema"`
	Variables     []any              `json:"variables"`
	SamplePayload map[string]any     `json:"sample_payload"`
	OutputFormat  enums.OutputFormat `json:"output_format"`
	Checksum      *string            `json:"checksum"`
	IsPublished   bool               `json:"is_published"`
	PublishedAt   *time.Time         `json:"published_at"`
	CreatedBy     *string            `json:"created_by"`
	CreatedAt     time.Time          `json:"created_at"`
}

type CreateTemplateVersionRequest struct {
	Content       string             `json:"content"`
	Schema        map[string]any     `json:"schema"`
	Variables     []any              `json:"variables"`
	SamplePayload map[string]any     `json:"sample_payload"`
	OutputFormat  enums.OutputFormat `json:"output_format"`
	CreatedBy     *string            `json:"created_by"`
}

type TemplateVersionListResponse struct {
	Data []TemplateVersionResponse `json:"data"`
}

func VersionFromEntity(v verEntity.TemplateVersion, includeContent bool) TemplateVersionResponse {
	resp := TemplateVersionResponse{
		ID: v.ID, TenantID: v.TenantID, TemplateID: v.TemplateID, Version: v.Version,
		Schema: v.Schema, Variables: v.Variables, SamplePayload: v.SamplePayload,
		OutputFormat: v.OutputFormat, Checksum: v.Checksum, IsPublished: v.IsPublished,
		PublishedAt: v.PublishedAt, CreatedBy: v.CreatedBy, CreatedAt: v.CreatedAt,
	}
	if includeContent {
		resp.Content = v.Content
	}
	return resp
}

func (r CreateTemplateVersionRequest) ToEntity(tenantID *string, templateID int64) verEntity.TemplateVersion {
	return verEntity.TemplateVersion{
		TenantID: tenantID, TemplateID: templateID, Content: r.Content,
		Schema: r.Schema, Variables: r.Variables, SamplePayload: r.SamplePayload,
		OutputFormat: r.OutputFormat, CreatedBy: r.CreatedBy,
	}
}
