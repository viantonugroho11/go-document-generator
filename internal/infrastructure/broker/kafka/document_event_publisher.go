package kafka

import (
	"context"

	docEntity "go-document-generator/internal/entity/documents"
	"go-document-generator/internal/transport/event/events"
	ucDoc "go-document-generator/internal/usecase/documents"

	libkafka "github.com/viantonugroho11/go-lib/kafka"
)

type DocumentEventPublisherKafka struct {
	queued    *libkafka.Producer[events.DocumentQueuedEvent]
	retried   *libkafka.Producer[events.DocumentRetriedEvent]
	generated *libkafka.Producer[events.DocumentGeneratedEvent]
	failed    *libkafka.Producer[events.DocumentFailedEvent]
	cancelled *libkafka.Producer[events.DocumentCancelledEvent]
}

func NewDocumentEventPublisherKafka(
	queued *libkafka.Producer[events.DocumentQueuedEvent],
	retried *libkafka.Producer[events.DocumentRetriedEvent],
	generated *libkafka.Producer[events.DocumentGeneratedEvent],
	failed *libkafka.Producer[events.DocumentFailedEvent],
	cancelled *libkafka.Producer[events.DocumentCancelledEvent],
) ucDoc.DocumentEventPublisher {
	return &DocumentEventPublisherKafka{
		queued:    queued,
		retried:   retried,
		generated: generated,
		failed:    failed,
		cancelled: cancelled,
	}
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

func (p *DocumentEventPublisherKafka) PublishDocumentGenerated(ctx context.Context, d docEntity.Document) error {
	return p.generated.Publish(ctx, events.DocumentGeneratedEvent{
		ID: d.ID, RequestID: d.RequestID, Status: string(d.Status),
		OutputFormat: string(d.OutputFormat), FilePath: d.FilePath, FileSize: d.FileSize,
	})
}

func (p *DocumentEventPublisherKafka) PublishDocumentFailed(ctx context.Context, d docEntity.Document) error {
	return p.failed.Publish(ctx, events.DocumentFailedEvent{
		ID: d.ID, RequestID: d.RequestID, Status: string(d.Status), ErrorMessage: d.ErrorMessage,
	})
}

func (p *DocumentEventPublisherKafka) PublishDocumentCancelled(ctx context.Context, d docEntity.Document) error {
	return p.cancelled.Publish(ctx, events.DocumentCancelledEvent{
		ID: d.ID, RequestID: d.RequestID, Status: string(d.Status),
	})
}
