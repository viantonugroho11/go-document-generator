package states

import (
	"context"

	docEntity "go-document-generator/internal/entity/documents"
	"go-document-generator/internal/entity/enums"
	"go-document-generator/internal/shared/apperror"
)

type pending struct {
	stateMachine *documentStateMachine
	h            Handlers
}

func (s pending) Do(ctx context.Context, update docEntity.Document) (docEntity.Document, error) {
	s.stateMachine.data = &update

	switch update.Status {
	case enums.DocumentStatusQueued:
		return s.h.OnToQueued.OnStateTransition(ctx, update)
	case enums.DocumentStatusCancelled:
		return s.h.OnToCancelled.OnStateTransition(ctx, update)
	case enums.DocumentStatusPending:
		return s.h.OnFieldUpdate.OnStateTransition(ctx, update)
	default:
		return docEntity.Document{}, apperror.ErrInvalidState
	}
}
