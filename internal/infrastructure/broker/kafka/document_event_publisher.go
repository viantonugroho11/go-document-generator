package kafka

import (
	"context"

	docEntity "go-document-generator/internal/entity/documents"
	"go-document-generator/internal/transport/event/events"
	ucDoc "go-document-generator/internal/usecase/documents"

	libkafka "github.com/viantonugroho11/go-lib/kafka"
)

type DocumentEventPublisherKafka struct {
	queued  *libkafka.Producer[events.DocumentQueuedEvent]
	retried *libkafka.Producer[events.DocumentRetriedEvent]
}

func NewDocumentEventPublisherKafka(
	queued *libkafka.Producer[events.DocumentQueuedEvent],
	retried *libkafka.Producer[events.DocumentRetriedEvent],
) ucDoc.DocumentEventPublisher {
	return &DocumentEventPublisherKafka{queued: queued, retried: retried}
}

func (p *DocumentEventPublisherKafka) PublishDocumentQueued(ctx context.Context, d docEntity.Document) error {
	return p.queued.Publish(ctx, events.DocumentQueuedEvent{
		ID: d.ID, RequestID: d.RequestID, TemplateCode: d.TemplateCode, Status: string(d.Status),
	})
}

func (p *DocumentEventPublisherKafka) PublishDocumentRetried(ctx context.Context, d docEntity.Document) error {
	return p.retried.Publish(ctx, events.DocumentRetriedEvent{
		ID: d.ID, RequestID: d.RequestID, Status: string(d.Status),
	})
}
