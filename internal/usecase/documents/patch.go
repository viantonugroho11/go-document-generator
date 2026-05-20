package documents

import (
	"context"

	docEntity "go-document-generator/internal/entity/documents"
	"go-document-generator/internal/entity/enums"
)

func (s *service) applyStateMachine(ctx context.Context, existing docEntity.Document, update docEntity.Document) (docEntity.Document, error) {
	sm, err := s.stateMachine.NewStateMachine(&existing)
	if err != nil {
		return docEntity.Document{}, err
	}

	result, err := sm.Do(ctx, update)
	if err != nil {
		if result.Status == enums.DocumentStatusFailed {
			if saved, updErr := s.docs.Update(ctx, nil, result); updErr == nil {
				return saved, err
			}
		}
		return docEntity.Document{}, err
	}
	return result, nil
}

func (s *service) transitionDocument(ctx context.Context, existing docEntity.Document, target enums.DocumentStatus) (docEntity.Document, error) {
	update := existing
	update.Status = target
	return s.applyStateMachine(ctx, existing, update)
}
