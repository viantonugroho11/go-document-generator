package documenttemplates

import (
	"context"

	tplEntity "go-document-generator/internal/entity/documenttemplates"
)

type TemplateEventPublisher interface {
	PublishTemplateCreated(ctx context.Context, t tplEntity.Template) error
	PublishTemplateUpdated(ctx context.Context, t tplEntity.Template) error
}

type noopTemplatePublisher struct{}

func (noopTemplatePublisher) PublishTemplateCreated(context.Context, tplEntity.Template) error { return nil }
func (noopTemplatePublisher) PublishTemplateUpdated(context.Context, tplEntity.Template) error { return nil }

func NoopTemplatePublisher() TemplateEventPublisher { return noopTemplatePublisher{} }
