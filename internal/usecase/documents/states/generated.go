package states

import (
	"context"

	docEntity "go-document-generator/internal/entity/documents"
	"go-document-generator/internal/entity/enums"
)

type generated struct {
	stateMachine *documentStateMachine
	h            Handlers
}

func (s generated) Do(ctx context.Context, update docEntity.Document) (docEntity.Document, error) {
	s.stateMachine.data = &update
	if update.Status == enums.DocumentStatusGenerated {
		return update, nil
	}
	return s.h.OnTerminal.OnStateTransition(ctx, update)
}
