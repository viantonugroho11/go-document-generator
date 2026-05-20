package documenttemplateversions

import (
	"context"

	verEntity "go-document-generator/internal/entity/documenttemplateversions"
)

type VersionEventPublisher interface {
	PublishVersionCreated(ctx context.Context, v verEntity.TemplateVersion) error
	PublishVersionPublished(ctx context.Context, v verEntity.TemplateVersion) error
}

type noopVersionPublisher struct{}

func (noopVersionPublisher) PublishVersionCreated(context.Context, verEntity.TemplateVersion) error   { return nil }
func (noopVersionPublisher) PublishVersionPublished(context.Context, verEntity.TemplateVersion) error { return nil }

func NoopVersionPublisher() VersionEventPublisher { return noopVersionPublisher{} }
