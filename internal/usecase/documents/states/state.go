package states

import (
	"context"
	"fmt"

	docEntity "go-document-generator/internal/entity/documents"
	"go-document-generator/internal/entity/enums"
)

// IOnStateTransition dieksekusi saat dokumen berpindah ke status target.
type IOnStateTransition interface {
	OnStateTransition(ctx context.Context, update docEntity.Document) (docEntity.Document, error)
}

// Handlers handler transisi yang diinjeksikan ke state machine.
type Handlers struct {
	OnFieldUpdate  IOnStateTransition
	OnToQueued     IOnStateTransition
	OnToProcessing IOnStateTransition
	OnToGenerated  IOnStateTransition
	OnToCancelled  IOnStateTransition
	OnToFailed     IOnStateTransition
	OnRetry        IOnStateTransition
	OnTerminal     IOnStateTransition
}

type IDocumentState interface {
	Do(ctx context.Context, update docEntity.Document) (docEntity.Document, error)
}

type IDocumentStateMachine interface {
	IDocumentState
	Document() *docEntity.Document
}

type IDocumentStateMachineFactory interface {
	NewStateMachine(current *docEntity.Document) (IDocumentStateMachine, error)
}

type documentStateMachine struct {
	data    *docEntity.Document
	current IDocumentState
	h       Handlers

	pending    IDocumentState
	queued     IDocumentState
	processing IDocumentState
	generated  IDocumentState
	failed     IDocumentState
	cancelled  IDocumentState
}

type documentStateMachineFactory struct {
	handlers Handlers
}

func NewDocumentStateMachineFactory(handlers Handlers) IDocumentStateMachineFactory {
	return &documentStateMachineFactory{handlers: handlers}
}

func (f *documentStateMachineFactory) NewStateMachine(current *docEntity.Document) (IDocumentStateMachine, error) {
	if current == nil || current.ID <= 0 {
		return nil, fmt.Errorf("document ID is required")
	}

	sm := &documentStateMachine{
		data: current,
		h:    f.handlers,
	}

	sm.pending = pending{stateMachine: sm, h: f.handlers}
	sm.queued = queued{stateMachine: sm, h: f.handlers}
	sm.processing = processing{stateMachine: sm, h: f.handlers}
	sm.generated = generated{stateMachine: sm, h: f.handlers}
	sm.failed = failed{stateMachine: sm, h: f.handlers}
	sm.cancelled = cancelled{stateMachine: sm, h: f.handlers}

	switch current.Status {
	case enums.DocumentStatusPending:
		sm.current = sm.pending
	case enums.DocumentStatusQueued:
		sm.current = sm.queued
	case enums.DocumentStatusProcessing:
		sm.current = sm.processing
	case enums.DocumentStatusGenerated:
		sm.current = sm.generated
	case enums.DocumentStatusFailed:
		sm.current = sm.failed
	case enums.DocumentStatusCancelled:
		sm.current = sm.cancelled
	default:
		return nil, fmt.Errorf("unknown document status: %s", current.Status)
	}

	return sm, nil
}

func (s *documentStateMachine) Do(ctx context.Context, update docEntity.Document) (docEntity.Document, error) {
	return s.current.Do(ctx, update)
}

func (s *documentStateMachine) Document() *docEntity.Document {
	return s.data
}
