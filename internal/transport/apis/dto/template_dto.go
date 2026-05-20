package dto

import (
	"time"

	tplEntity "go-document-generator/internal/entity/documenttemplates"
	"go-document-generator/internal/entity/enums"
	"go-document-generator/internal/shared/tenant"
)

type DocumentTemplateResponse struct {
	ID            int64              `json:"id"`
	TenantID      *string            `json:"tenant_id"`
	Code          string             `json:"code"`
	Name          string             `json:"name"`
	Description   *string            `json:"description"`
	Engine        enums.TemplateEngine `json:"engine"`
	DefaultFormat enums.OutputFormat `json:"default_format"`
	Category      *string            `json:"category"`
	IsActive      bool               `json:"is_active"`
	CreatedBy     *string            `json:"created_by"`
	UpdatedBy     *string            `json:"updated_by"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}

type CreateTemplateRequest struct {
	TenantID      *string            `json:"tenant_id"`
	Code          string             `json:"code"`
	Name          string             `json:"name"`
	Description   *string            `json:"description"`
	Engine        enums.TemplateEngine `json:"engine"`
	DefaultFormat enums.OutputFormat `json:"default_format"`
	Category      *string            `json:"category"`
	IsActive      *bool              `json:"is_active"`
	CreatedBy     *string            `json:"created_by"`
}

type PatchTemplateRequest struct {
	Name          *string             `json:"name"`
	Description   *string             `json:"description"`
	Engine        *enums.TemplateEngine `json:"engine"`
	DefaultFormat *enums.OutputFormat `json:"default_format"`
	Category      *string             `json:"category"`
	IsActive      *bool               `json:"is_active"`
	UpdatedBy     *string             `json:"updated_by"`
}

type TemplateListResponse struct {
	Data []DocumentTemplateResponse `json:"data"`
	Meta PaginationMeta             `json:"meta"`
}

func TemplateFromEntity(t tplEntity.Template) DocumentTemplateResponse {
	return DocumentTemplateResponse{
		ID: t.ID, TenantID: t.TenantID, Code: t.Code, Name: t.Name, Description: t.Description,
		Engine: t.Engine, DefaultFormat: t.DefaultFormat, Category: t.Category, IsActive: t.IsActive,
		CreatedBy: t.CreatedBy, UpdatedBy: t.UpdatedBy, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
	}
}

func (r CreateTemplateRequest) ToEntity(headerTenant *string) tplEntity.Template {
	tid := r.TenantID
	if tid == nil {
		tid = headerTenant
	}
	active := true
	if r.IsActive != nil {
		active = *r.IsActive
	}
	return tplEntity.Template{
		TenantID: tid, Code: r.Code, Name: r.Name, Description: r.Description,
		Engine: r.Engine, DefaultFormat: r.DefaultFormat, Category: r.Category,
		IsActive: active, CreatedBy: r.CreatedBy,
	}
}

func ApplyPatchTemplate(existing tplEntity.Template, req PatchTemplateRequest) tplEntity.Template {
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Description != nil {
		existing.Description = req.Description
	}
	if req.Engine != nil {
		existing.Engine = *req.Engine
	}
	if req.DefaultFormat != nil {
		existing.DefaultFormat = *req.DefaultFormat
	}
	if req.Category != nil {
		existing.Category = req.Category
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	if req.UpdatedBy != nil {
		existing.UpdatedBy = req.UpdatedBy
	}
	return existing
}

// ResolveTenant menggabungkan header dan body tenant_id.
func ResolveTenant(header, body *string) *string {
	if body != nil && *body != "" {
		return tenant.PtrOrNil(*body)
	}
	return header
}
