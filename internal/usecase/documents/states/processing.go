package states

import (
	"context"

	docEntity "go-document-generator/internal/entity/documents"
	"go-document-generator/internal/entity/enums"
	"go-document-generator/internal/shared/apperror"
)

type processing struct {
	stateMachine *documentStateMachine
	h            Handlers
}

func (s processing) Do(ctx context.Context, update docEntity.Document) (docEntity.Document, error) {
	s.stateMachine.data = &update

	switch update.Status {
	case enums.DocumentStatusGenerated:
		return s.h.OnToGenerated.OnStateTransition(ctx, update)
	case enums.DocumentStatusFailed:
		return s.h.OnToFailed.OnStateTransition(ctx, update)
	case enums.DocumentStatusCancelled:
		return s.h.OnToCancelled.OnStateTransition(ctx, update)
	case enums.DocumentStatusProcessing:
		return update, nil
	default:
		return docEntity.Document{}, apperror.ErrInvalidState
	}
}
