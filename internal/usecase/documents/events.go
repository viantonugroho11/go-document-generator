package documents

import (
	"context"

	docEntity "go-document-generator/internal/entity/documents"
)

type DocumentEventPublisher interface {
	// Observability events (document-events topic)
	PublishDocumentQueued(ctx context.Context, d docEntity.Document) error
	PublishDocumentRetried(ctx context.Context, d docEntity.Document) error
	PublishDocumentGenerated(ctx context.Context, d docEntity.Document) error
	PublishDocumentFailed(ctx context.Context, d docEntity.Document) error
	PublishDocumentCancelled(ctx context.Context, d docEntity.Document) error
	PublishDocumentsZipped(ctx context.Context, ids []int64, tenantID *string, zipPath, outputFormat string) error
	PublishDocumentsMerged(ctx context.Context, ids []int64, tenantID *string, mergedPath, outputFormat string) error

	// Processing trigger (document-process topic)
	// Dipanggil setiap kali dokumen perlu dirender: saat Create dan saat Retry.
	PublishDocumentProcess(ctx context.Context, d docEntity.Document) error
}

type noopDocumentPublisher struct{}

func (noopDocumentPublisher) PublishDocumentQueued(context.Context, docEntity.Document) error    { return nil }
func (noopDocumentPublisher) PublishDocumentRetried(context.Context, docEntity.Document) error   { return nil }
func (noopDocumentPublisher) PublishDocumentGenerated(context.Context, docEntity.Document) error { return nil }
func (noopDocumentPublisher) PublishDocumentFailed(context.Context, docEntity.Document) error    { return nil }
func (noopDocumentPublisher) PublishDocumentCancelled(context.Context, docEntity.Document) error { return nil }
func (noopDocumentPublisher) PublishDocumentsZipped(context.Context, []int64, *string, string, string) error { return nil }
func (noopDocumentPublisher) PublishDocumentsMerged(context.Context, []int64, *string, string, string) error { return nil }
func (noopDocumentPublisher) PublishDocumentProcess(context.Context, docEntity.Document) error  { return nil }

func NoopDocumentPublisher() DocumentEventPublisher { return noopDocumentPublisher{} }
