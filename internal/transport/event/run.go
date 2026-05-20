package event

import (
	"context"

	"go-document-generator/internal/config"
	infrakafka "go-document-generator/internal/infrastructure/broker/kafka"
	transportkafka "go-document-generator/internal/transport/event/kafka"
	usecaseusers "go-document-generator/internal/usecase/users"
)

const (
	ConsumerNameUser  = "user"
	ConsumerNameOrder = "order"
)

// ConsumerNames daftar nama consumer yang didukung (flag -consumer).
func ConsumerNames() []string {
	return []string{ConsumerNameUser, ConsumerNameOrder}
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
