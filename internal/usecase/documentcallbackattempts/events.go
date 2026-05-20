package documentcallbackattempts

import (
	"context"

	cbEntity "go-document-generator/internal/entity/documentcallbackattempts"
)

type CallbackEventPublisher interface {
	PublishCallbackAttempt(ctx context.Context, a cbEntity.CallbackAttempt) error
}

type noopCallbackPublisher struct{}

func (noopCallbackPublisher) PublishCallbackAttempt(context.Context, cbEntity.CallbackAttempt) error {
	return nil
}

func NoopCallbackPublisher() CallbackEventPublisher { return noopCallbackPublisher{} }
