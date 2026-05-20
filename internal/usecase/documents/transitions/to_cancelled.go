package transitions

import (
	"context"

	docEntity "go-document-generator/internal/entity/documents"
	"go-document-generator/internal/entity/enums"
)

type toCancelled struct{}

func NewToCancelled() *toCancelled { return &toCancelled{} }

func (h *toCancelled) OnStateTransition(_ context.Context, update docEntity.Document) (docEntity.Document, error) {
	update.Status = enums.DocumentStatusCancelled
	return update, nil
}
