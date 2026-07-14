package kafka

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	docEntity "go-document-generator/internal/entity/documents"
	"go-document-generator/internal/transport/event/events"
	ucDoc "go-document-generator/internal/usecase/documents"

	libkafka "github.com/viantonugroho11/go-lib/kafka"
)

// DocumentEventPublisherKafka mempublikasikan event dokumen ke Kafka.
// doc    → topic document-events  (lifecycle events: CREATE / UPDATE)
// bulk   → topic document-events  (zip / merge completion events)
// process → topic document-process (generation trigger)
type DocumentEventPublisherKafka struct {
	doc     *libkafka.Producer[events.DocumentEvent]     // document-events
	bulk    *libkafka.Producer[events.DocumentBulkEvent] // document-events (bulk ops)
	process *libkafka.Producer[events.DocumentEvent]     // document-process
}

// NewDocumentEventPublisherKafka membuat publisher dengan 3 producer.
func NewDocumentEventPublisherKafka(
	doc *libkafka.Producer[events.DocumentEvent],
	bulk *libkafka.Producer[events.DocumentBulkEvent],
	process *libkafka.Producer[events.DocumentEvent],
) ucDoc.DocumentEventPublisher {
	return &DocumentEventPublisherKafka{doc: doc, bulk: bulk, process: process}
}

func (p *DocumentEventPublisherKafka) PublishDocumentEvent(ctx context.Context, action string, before, after *docEntity.Document) error {
	resourceID := ""
	if after != nil {
		resourceID = after.RequestID
	}
	evt := events.DocumentEvent{
		ResourceID: resourceID,
		Meta: events.EventMeta{
			EventID:              newEventID(),
			EventTimestamp:       time.Now().UTC().Format(time.RFC3339),
			Action:               action,
			Resource:             "Document",
			MessageSchemaVersion: events.MessageSchemaVersion,
		},
		Before: toDocumentState(before),
		After:  toDocumentState(after),
	}
	return p.doc.Publish(ctx, evt)
}

func (p *DocumentEventPublisherKafka) PublishDocumentBulkEvent(
	ctx context.Context,
	resource string,
	ids []int64,
	tenantID *string,
	outputPath, outputFormat string,
) error {
	after := &events.DocumentBulkState{
		DocumentIDs:  ids,
		TenantID:     tenantID,
		OutputPath:   outputPath,
		OutputFormat: outputFormat,
	}
	evt := events.DocumentBulkEvent{
		ResourceID: outputPath,
		Meta: events.EventMeta{
			EventID:              newEventID(),
			EventTimestamp:       time.Now().UTC().Format(time.RFC3339),
			Action:               "CREATE",
			Resource:             resource,
			MessageSchemaVersion: events.MessageSchemaVersion,
		},
		Before: nil,
		After:  after,
	}
	return p.bulk.Publish(ctx, evt)
}

func (p *DocumentEventPublisherKafka) PublishDocumentProcess(ctx context.Context, d docEntity.Document) error {
	evt := events.DocumentEvent{
		ResourceID: d.RequestID,
		Meta: events.EventMeta{
			EventID:              newEventID(),
			EventTimestamp:       time.Now().UTC().Format(time.RFC3339),
			Action:               "CREATE",
			Resource:             "Document",
			MessageSchemaVersion: events.MessageSchemaVersion,
		},
		Before: nil,
		After:  toDocumentState(&d),
	}
	return p.process.Publish(ctx, evt)
}

// toDocumentState mengkonversi entity dokumen ke snapshot untuk event envelope.
func toDocumentState(d *docEntity.Document) *events.DocumentState {
	if d == nil {
		return nil
	}
	var sp *string
	if d.StorageProvider != nil {
		s := string(*d.StorageProvider)
		sp = &s
	}
	return &events.DocumentState{
		ID:              d.ID,
		RequestID:       d.RequestID,
		TenantID:        d.TenantID,
		TemplateCode:    d.TemplateCode,
		TemplateVersion: d.TemplateVersion,
		OutputFormat:    string(d.OutputFormat),
		Status:          string(d.Status),
		ErrorMessage:    d.ErrorMessage,
		RetryCount:      d.RetryCount,
		FilePath:        d.FilePath,
		FileName:        d.FileName,
		FileSize:        d.FileSize,
		ContentType:     d.ContentType,
		Checksum:        d.Checksum,
		StorageProvider: sp,
		ProcessedAt:     d.ProcessedAt,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}
}

// newEventID membuat UUID v4 sederhana tanpa dependency external.
func newEventID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
