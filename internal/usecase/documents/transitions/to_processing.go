package transitions

import (
	"context"

	docEntity "go-document-generator/internal/entity/documents"
	"go-document-generator/internal/entity/enums"
)

type toProcessing struct{}

func NewToProcessing() *toProcessing { return &toProcessing{} }

func (h *toProcessing) OnStateTransition(_ context.Context, update docEntity.Document) (docEntity.Document, error) {
	update.Status = enums.DocumentStatusProcessing
	return update, nil
}
