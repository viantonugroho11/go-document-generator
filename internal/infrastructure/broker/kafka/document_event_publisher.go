package kafka

import (
	"context"
	"time"

	docEntity "go-document-generator/internal/entity/documents"
	"go-document-generator/internal/transport/event/events"
	ucDoc "go-document-generator/internal/usecase/documents"

	libkafka "github.com/viantonugroho11/go-lib/kafka"
)

type DocumentEventPublisherKafka struct {
	queued    *libkafka.Producer[events.DocumentQueuedEvent]    // document-events (observability)
	retried   *libkafka.Producer[events.DocumentRetriedEvent]
	generated *libkafka.Producer[events.DocumentGeneratedEvent]
	failed    *libkafka.Producer[events.DocumentFailedEvent]
	cancelled *libkafka.Producer[events.DocumentCancelledEvent]
	zipped    *libkafka.Producer[events.DocumentsZippedEvent]
	merged    *libkafka.Producer[events.DocumentsMergedEvent]
	process   *libkafka.Producer[events.DocumentQueuedEvent]   // document-process (generation trigger)
}

func NewDocumentEventPublisherKafka(
	queued *libkafka.Producer[events.DocumentQueuedEvent],
	retried *libkafka.Producer[events.DocumentRetriedEvent],
	generated *libkafka.Producer[events.DocumentGeneratedEvent],
	failed *libkafka.Producer[events.DocumentFailedEvent],
	cancelled *libkafka.Producer[events.DocumentCancelledEvent],
	zipped *libkafka.Producer[events.DocumentsZippedEvent],
	merged *libkafka.Producer[events.DocumentsMergedEvent],
	process *libkafka.Producer[events.DocumentQueuedEvent],
) ucDoc.DocumentEventPublisher {
	return &DocumentEventPublisherKafka{
		queued:    queued,
		retried:   retried,
		generated: generated,
		failed:    failed,
		cancelled: cancelled,
		zipped:    zipped,
		merged:    merged,
		process:   process,
	}
}

func (p *DocumentEventPublisherKafka) PublishDocumentQueued(ctx context.Context, d docEntity.Document) error {
	return p.queued.Publish(ctx, events.DocumentQueuedEvent{
		ID: d.ID, RequestID: d.RequestID, TenantID: d.TenantID,
		TemplateCode: d.TemplateCode, TemplateVersion: d.TemplateVersion,
		OutputFormat: string(d.OutputFormat), Status: string(d.Status),
		OccurredAt: time.Now().UTC(),
	})
}

func (p *DocumentEventPublisherKafka) PublishDocumentProcess(ctx context.Context, d docEntity.Document) error {
	return p.process.Publish(ctx, events.DocumentQueuedEvent{
		ID: d.ID, RequestID: d.RequestID, TenantID: d.TenantID,
		TemplateCode: d.TemplateCode, TemplateVersion: d.TemplateVersion,
		OutputFormat: string(d.OutputFormat), Status: string(d.Status),
		OccurredAt: time.Now().UTC(),
	})
}

func (p *DocumentEventPublisherKafka) PublishDocumentRetried(ctx context.Context, d docEntity.Document) error {
	return p.retried.Publish(ctx, events.DocumentRetriedEvent{
		ID: d.ID, RequestID: d.RequestID, TenantID: d.TenantID,
		RetryCount: d.RetryCount, Status: string(d.Status),
		OccurredAt: time.Now().UTC(),
	})
}

func (p *DocumentEventPublisherKafka) PublishDocumentGenerated(ctx context.Context, d docEntity.Document) error {
	var sp *string
	if d.StorageProvider != nil {
		s := string(*d.StorageProvider)
		sp = &s
	}
	return p.generated.Publish(ctx, events.DocumentGeneratedEvent{
		ID: d.ID, RequestID: d.RequestID, TenantID: d.TenantID,
		TemplateCode: d.TemplateCode, OutputFormat: string(d.OutputFormat), Status: string(d.Status),
		FileName: d.FileName, FilePath: d.FilePath, FileSize: d.FileSize,
		ContentType: d.ContentType, Checksum: d.Checksum, StorageProvider: sp,
		ProcessedAt: d.ProcessedAt, OccurredAt: time.Now().UTC(),
	})
}

func (p *DocumentEventPublisherKafka) PublishDocumentFailed(ctx context.Context, d docEntity.Document) error {
	return p.failed.Publish(ctx, events.DocumentFailedEvent{
		ID: d.ID, RequestID: d.RequestID, TenantID: d.TenantID,
		Status: string(d.Status), ErrorMessage: d.ErrorMessage,
		RetryCount: d.RetryCount, OccurredAt: time.Now().UTC(),
	})
}

func (p *DocumentEventPublisherKafka) PublishDocumentCancelled(ctx context.Context, d docEntity.Document) error {
	return p.cancelled.Publish(ctx, events.DocumentCancelledEvent{
		ID: d.ID, RequestID: d.RequestID, TenantID: d.TenantID,
		Status: string(d.Status), OccurredAt: time.Now().UTC(),
	})
}

func (p *DocumentEventPublisherKafka) PublishDocumentsZipped(ctx context.Context, ids []int64, tenantID *string, zipPath, _ string) error {
	return p.zipped.Publish(ctx, events.DocumentsZippedEvent{
		DocumentIDs: ids, TenantID: tenantID,
		ZipPath: zipPath, OccurredAt: time.Now().UTC(),
	})
}

func (p *DocumentEventPublisherKafka) PublishDocumentsMerged(ctx context.Context, ids []int64, tenantID *string, mergedPath, outputFormat string) error {
	return p.merged.Publish(ctx, events.DocumentsMergedEvent{
		DocumentIDs: ids, TenantID: tenantID,
		MergedPath: mergedPath, OutputFormat: outputFormat, OccurredAt: time.Now().UTC(),
	})
}
