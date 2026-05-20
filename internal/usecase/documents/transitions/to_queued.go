package transitions

import (
	"context"

	docEntity "go-document-generator/internal/entity/documents"
	"go-document-generator/internal/entity/enums"
)

type toQueued struct{}

func NewToQueued() *toQueued { return &toQueued{} }

func (h *toQueued) OnStateTransition(_ context.Context, update docEntity.Document) (docEntity.Document, error) {
	update.Status = enums.DocumentStatusQueued
	update.ErrorMessage = nil
	return update, nil
}
