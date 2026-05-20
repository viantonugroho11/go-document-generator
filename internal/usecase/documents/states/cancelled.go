package states

import (
	"context"

	docEntity "go-document-generator/internal/entity/documents"
	"go-document-generator/internal/entity/enums"
)

type cancelled struct {
	stateMachine *documentStateMachine
	h            Handlers
}

func (s cancelled) Do(ctx context.Context, update docEntity.Document) (docEntity.Document, error) {
	s.stateMachine.data = &update
	if update.Status == enums.DocumentStatusCancelled {
		return update, nil
	}
	return s.h.OnTerminal.OnStateTransition(ctx, update)
}
