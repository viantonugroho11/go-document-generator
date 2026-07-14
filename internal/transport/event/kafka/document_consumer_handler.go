package kafka

import (
	"context"
	"log"

	"go-document-generator/internal/transport/event/events"
	ucDoc "go-document-generator/internal/usecase/documents"

	libkafka "github.com/viantonugroho11/go-lib/kafka"
)

// DocumentProcessHandler consume DocumentQueuedEvent dari topic document-process
// dan menjalankan pipeline generation: QUEUED → PROCESSING → GENERATED/FAILED.
type DocumentProcessHandler struct {
	docs ucDoc.Service
}

func NewDocumentProcessHandler(docs ucDoc.Service) *DocumentProcessHandler {
	return &DocumentProcessHandler{docs: docs}
}

func (h *DocumentProcessHandler) Name() string { return "document-process" }

func (h *DocumentProcessHandler) Handle(ctx context.Context, evt events.DocumentQueuedEvent, _ ...libkafka.Header) libkafka.Progress {
	if evt.ID == 0 {
		return libkafka.Progress{Status: libkafka.ProgressDrop, Result: "document id is zero"}
	}

	if err := h.docs.Process(ctx, evt.ID, evt.TenantID); err != nil {
		log.Printf("document_consumer: Process id=%d request_id=%s: %v", evt.ID, evt.RequestID, err)
		return libkafka.Progress{Status: libkafka.ProgressError, Result: err.Error()}
	}

	log.Printf("document_consumer: generated id=%d request_id=%s", evt.ID, evt.RequestID)
	return libkafka.Progress{Status: libkafka.ProgressSuccess, Result: "generated"}
}
