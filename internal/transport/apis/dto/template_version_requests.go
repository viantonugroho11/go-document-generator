package dto

import entityVer "go-document-generator/internal/entity/documenttemplateversions"

type CreateTemplateVersionRequest struct {
	Content       string                 `json:"content"`
	Schema        map[string]any         `json:"schema"`
	SamplePayload map[string]any         `json:"sample_payload"`
	IsPublished   *bool                  `json:"is_published,omitempty"`
	Version       int                    `json:"version,omitempty"` // optional; auto if 0
}

func (r *CreateTemplateVersionRequest) ToEntity() entityVer.DocumentTemplateVersion {
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

type UpdateTemplateVersionRequest struct {
	Content       string                 `json:"content"`
	Schema        map[string]any         `json:"schema"`
	SamplePayload map[string]any         `json:"sample_payload"`
	IsPublished   *bool                  `json:"is_published,omitempty"`
}

