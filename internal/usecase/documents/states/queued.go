package states

import (
	"context"

	docEntity "go-document-generator/internal/entity/documents"
	"go-document-generator/internal/entity/enums"
	"go-document-generator/internal/shared/apperror"
)

type queued struct {
	stateMachine *documentStateMachine
	h            Handlers
}

func (s queued) Do(ctx context.Context, update docEntity.Document) (docEntity.Document, error) {
	s.stateMachine.data = &update

	switch update.Status {
	case enums.DocumentStatusProcessing:
		return s.h.OnToProcessing.OnStateTransition(ctx, update)
	case enums.DocumentStatusCancelled:
		return s.h.OnToCancelled.OnStateTransition(ctx, update)
	case enums.DocumentStatusQueued:
		return s.h.OnFieldUpdate.OnStateTransition(ctx, update)
	default:
		return docEntity.Document{}, apperror.ErrInvalidState
	}
}
