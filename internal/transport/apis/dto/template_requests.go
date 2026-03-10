package dto

import (
	entityTmpl "go-document-generator/internal/entity/documenttemplates"
	entityVer "go-document-generator/internal/entity/documenttemplateversions"
)

type CreateOrAddTemplateWithVersionRequest struct {
	Code         string                 `json:"code"`
	Name         string                 `json:"name"`
	Description  *string                `json:"description"`
	Engine       string                 `json:"engine"`
	OutputFormat string                 `json:"output_format"`

	// initial/new version
	Content       string                 `json:"content"`
	Schema        map[string]any         `json:"schema"`
	SamplePayload map[string]any         `json:"sample_payload"`
	IsPublished   *bool                  `json:"is_published,omitempty"`
	Version       int                    `json:"version,omitempty"` // optional; auto if 0
}

func (r *CreateOrAddTemplateWithVersionRequest) ToTemplateEntity() entityTmpl.DocumentTemplate {
	return entityTmpl.DocumentTemplate{
		Code:         r.Code,
		Name:         r.Name,
		Description:  r.Description,
		Engine:       r.Engine,
		OutputFormat: r.OutputFormat,
		// IsActive akan diset di usecase saat create awal
	}
}

func (r *CreateOrAddTemplateWithVersionRequest) ToVersionEntity() entityVer.DocumentTemplateVersion {
	ver := entityVer.DocumentTemplateVersion{
		Content:       r.Content,
		Schema:        r.Schema,
		SamplePayload: r.SamplePayload,
		Version:       r.Version,
	}
	if r.IsPublished != nil {
		ver.IsPublished = *r.IsPublished
	}
	return ver
}

type UpdateTemplateRequest struct {
	Name         string  `json:"name"`
	Description  *string `json:"description"`
	Engine       string  `json:"engine"`
	OutputFormat string  `json:"output_format"`
	IsActive     *bool   `json:"is_active"`
}

