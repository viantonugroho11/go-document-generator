package states

import (
	"context"

	docEntity "go-document-generator/internal/entity/documents"
	"go-document-generator/internal/entity/enums"
	"go-document-generator/internal/shared/apperror"
)

type failed struct {
	stateMachine *documentStateMachine
	h            Handlers
}

func (s failed) Do(ctx context.Context, update docEntity.Document) (docEntity.Document, error) {
	s.stateMachine.data = &update

	switch update.Status {
	case enums.DocumentStatusQueued:
		return s.h.OnRetry.OnStateTransition(ctx, update)
	default:
		return docEntity.Document{}, apperror.ErrInvalidState
	}
}
