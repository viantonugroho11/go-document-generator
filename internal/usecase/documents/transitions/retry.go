package transitions

import (
	"context"

	docEntity "go-document-generator/internal/entity/documents"
	"go-document-generator/internal/entity/enums"
)

type retry struct{}

func NewRetry() *retry { return &retry{} }

func (h *retry) OnStateTransition(_ context.Context, update docEntity.Document) (docEntity.Document, error) {
	update.Status = enums.DocumentStatusQueued
	update.RetryCount++
	update.ErrorMessage = nil
	return update, nil
}
