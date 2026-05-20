package documents

import (
	"context"

	docEntity "go-document-generator/internal/entity/documents"
)

type DocumentEventPublisher interface {
	PublishDocumentQueued(ctx context.Context, d docEntity.Document) error
	PublishDocumentRetried(ctx context.Context, d docEntity.Document) error
}

type noopDocumentPublisher struct{}

func (noopDocumentPublisher) PublishDocumentQueued(context.Context, docEntity.Document) error  { return nil }
func (noopDocumentPublisher) PublishDocumentRetried(context.Context, docEntity.Document) error { return nil }

func NoopDocumentPublisher() DocumentEventPublisher { return noopDocumentPublisher{} }
