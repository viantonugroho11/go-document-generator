package transitions

import (
	"context"

	docEntity "go-document-generator/internal/entity/documents"
	"go-document-generator/internal/entity/enums"
)

type toFailed struct{}

func NewToFailed() *toFailed { return &toFailed{} }

func (h *toFailed) OnStateTransition(_ context.Context, update docEntity.Document) (docEntity.Document, error) {
	update.Status = enums.DocumentStatusFailed
	return update, nil
}
