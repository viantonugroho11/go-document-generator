package states

import (
	"context"
	"testing"

	docEntity "go-document-generator/internal/entity/documents"
	"go-document-generator/internal/entity/enums"
)

type stubTransition struct {
	apply func(docEntity.Document) docEntity.Document
}

func (s stubTransition) OnStateTransition(_ context.Context, update docEntity.Document) (docEntity.Document, error) {
	if s.apply != nil {
		return s.apply(update), nil
	}
	return update, nil
}

func testHandlers() Handlers {
	return Handlers{
		OnFieldUpdate:  stubTransition{},
		OnToQueued:     stubTransition{apply: func(d docEntity.Document) docEntity.Document { d.Status = enums.DocumentStatusQueued; return d }},
		OnToProcessing: stubTransition{apply: func(d docEntity.Document) docEntity.Document { d.Status = enums.DocumentStatusProcessing; return d }},
		OnToGenerated:  stubTransition{apply: func(d docEntity.Document) docEntity.Document { d.Status = enums.DocumentStatusGenerated; return d }},
		OnToCancelled:  stubTransition{apply: func(d docEntity.Document) docEntity.Document { d.Status = enums.DocumentStatusCancelled; return d }},
		OnToFailed:     stubTransition{apply: func(d docEntity.Document) docEntity.Document { d.Status = enums.DocumentStatusFailed; return d }},
		OnRetry:        stubTransition{apply: func(d docEntity.Document) docEntity.Document { d.Status = enums.DocumentStatusQueued; return d }},
		OnTerminal:     stubTransition{},
	}
}

func TestDocumentStateMachineTransitions(t *testing.T) {
	factory := NewDocumentStateMachineFactory(testHandlers())
	ctx := context.Background()

	doc := &docEntity.Document{ID: 1, Status: enums.DocumentStatusQueued, Payload: map[string]any{"x": 1}}

	sm, err := factory.NewStateMachine(doc)
	if err != nil {
		t.Fatal(err)
	}

	update := *doc
	update.Status = enums.DocumentStatusProcessing
	out, err := sm.Do(ctx, update)
	if err != nil || out.Status != enums.DocumentStatusProcessing {
		t.Fatalf("queued->processing: %+v err=%v", out, err)
	}

	sm, _ = factory.NewStateMachine(&out)
	update.Status = enums.DocumentStatusCancelled
	out, err = sm.Do(ctx, update)
	if err != nil || out.Status != enums.DocumentStatusCancelled {
		t.Fatalf("processing->cancelled: %+v err=%v", out, err)
	}
}
