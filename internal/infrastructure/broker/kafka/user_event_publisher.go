package kafka

import (
	"context"

	"go-document-generator/internal/entity/users"
	"go-document-generator/internal/transport/event/events"

	"github.com/viantonugroho11/go-lib/kafka"
)

// UserEventPublisherKafka implementasi users.UserEventPublisher (go-lib Producer).
type UserEventPublisherKafka struct {
	producer *kafka.Producer[events.UserCreatedEvent]
}

func NewUserEventPublisherKafka(producer *kafka.Producer[events.UserCreatedEvent]) *UserEventPublisherKafka {
	return &UserEventPublisherKafka{producer: producer}
}

func (p *UserEventPublisherKafka) PublishUser(ctx context.Context, user users.User) error {
	return p.producer.Publish(ctx, events.UserCreatedEvent{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}
