package documents

import (
	"context"

	docEntity "go-document-generator/internal/entity/documents"
)

// DocumentEventPublisher port untuk mempublikasikan event dokumen ke message broker.
type DocumentEventPublisher interface {
	// PublishDocumentEvent mengirim lifecycle event dokumen ke topic document-events.
	// action "CREATE" = before nil (dokumen baru).
	// action "UPDATE" = before berisi state sebelumnya.
	PublishDocumentEvent(ctx context.Context, action string, before, after *docEntity.Document) error

	// PublishDocumentBulkEvent mengirim event operasi zip / merge ke topic document-events.
	// resource = "DocumentZip" atau "DocumentMerge".
	PublishDocumentBulkEvent(ctx context.Context, resource string, ids []int64, tenantID *string, outputPath, outputFormat string) error

	// PublishDocumentProcess memicu generation worker via topic document-process.
	PublishDocumentProcess(ctx context.Context, d docEntity.Document) error
}

type noopDocumentPublisher struct{}

func (noopDocumentPublisher) PublishDocumentEvent(context.Context, string, *docEntity.Document, *docEntity.Document) error {
	return nil
}
func (noopDocumentPublisher) PublishDocumentBulkEvent(context.Context, string, []int64, *string, string, string) error {
	return nil
}
func (noopDocumentPublisher) PublishDocumentProcess(context.Context, docEntity.Document) error {
	return nil
}

func NoopDocumentPublisher() DocumentEventPublisher { return noopDocumentPublisher{} }
