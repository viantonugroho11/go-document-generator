package transitions

import (
	"context"

	docEntity "go-document-generator/internal/entity/documents"
	"go-document-generator/internal/shared/apperror"
	"go-document-generator/internal/shared/validators"
)

type fieldUpdate struct {
	deps Deps
}

func NewFieldUpdate(deps Deps) *fieldUpdate {
	return &fieldUpdate{deps: deps}
}

func (h *fieldUpdate) OnStateTransition(ctx context.Context, update docEntity.Document) (docEntity.Document, error) {
	if update.TemplateVersionID != nil && update.TemplateID != nil && len(update.Payload) > 0 {
		ver, err := h.deps.Versions.GetByID(ctx, nil, *update.TemplateID, *update.TemplateVersionID, update.TenantID)
		if err != nil {
			return docEntity.Document{}, err
		}
		if len(ver.Schema) > 0 {
			if err := validators.ValidateSchema(ver.Schema, update.Payload); err != nil {
				return docEntity.Document{}, err
			}
		}
	}
	return update, nil
}

// noopTransition tidak mengubah data (status terminal / no-op).
type noopTransition struct{}

func NewNoop() *noopTransition { return &noopTransition{} }

func (h *noopTransition) OnStateTransition(_ context.Context, update docEntity.Document) (docEntity.Document, error) {
	return update, apperror.ErrInvalidState
}
