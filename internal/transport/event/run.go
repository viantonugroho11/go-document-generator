package event

import (
	"context"

	"go-document-generator/internal/config"
	infrakafka "go-document-generator/internal/infrastructure/broker/kafka"
	transportkafka "go-document-generator/internal/transport/event/kafka"
	ucDoc "go-document-generator/internal/usecase/documents"
	usecaseusers "go-document-generator/internal/usecase/users"
)

const (
	ConsumerNameUser     = "user"
	ConsumerNameOrder    = "order"
	ConsumerNameDocument = "document"
)

// ConsumerNames daftar nama consumer yang didukung (flag -consumer).
func ConsumerNames() []string {
	return []string{ConsumerNameUser, ConsumerNameOrder, ConsumerNameDocument}
}

// RunUser menjalankan consumer Kafka untuk event user (topic & group dari cfg.Kafka).
func RunUser(ctx context.Context, cfg *config.Configuration, userService usecaseusers.UserService) (interface{ Close() error }, error) {
	h := transportkafka.NewUserCreatedHandler(userService)
	return infrakafka.RunWithConfig(ctx, cfg, cfg.Kafka.GroupID, cfg.Kafka.Topic, h)
}

// RunOrder menjalankan consumer Kafka untuk event order (topic_orders & group_id_orders).
func RunOrder(ctx context.Context, cfg *config.Configuration) (interface{ Close() error }, error) {
	h := transportkafka.NewOrderCreatedHandler()
	return infrakafka.RunWithConfig(ctx, cfg, cfg.Kafka.GroupIDOrders, cfg.Kafka.TopicOrders, h)
}

// RunDocument menjalankan consumer Kafka untuk document-process topic.
// Consumer ini memicu pipeline generation: QUEUED → PROCESSING → GENERATED/FAILED.
func RunDocument(ctx context.Context, cfg *config.Configuration, docs ucDoc.Service) (interface{ Close() error }, error) {
	topic := cfg.Kafka.TopicDocumentProcess
	if topic == "" {
		topic = "document-process"
	}
	groupID := cfg.Kafka.GroupIDDocumentWorker
	if groupID == "" {
		groupID = "document-generator-worker"
	}
	h := transportkafka.NewDocumentProcessHandler(docs)
	return infrakafka.RunWithConfig(ctx, cfg, groupID, topic, h)
}
