package kafka

import (
	"context"

	tplEntity "go-document-generator/internal/entity/documenttemplates"
	verEntity "go-document-generator/internal/entity/documenttemplateversions"
	"go-document-generator/internal/transport/event/events"
	ucTpl "go-document-generator/internal/usecase/documenttemplates"
	ucVer "go-document-generator/internal/usecase/documenttemplateversions"

	libkafka "github.com/viantonugroho11/go-lib/kafka"
)

type TemplateEventPublisherKafka struct {
	producer *libkafka.Producer[events.TemplateCreatedEvent]
}

func NewTemplateEventPublisherKafka(producer *libkafka.Producer[events.TemplateCreatedEvent]) ucTpl.TemplateEventPublisher {
	return &TemplateEventPublisherKafka{producer: producer}
}

func (p *TemplateEventPublisherKafka) PublishTemplateCreated(ctx context.Context, t tplEntity.Template) error {
	return p.producer.Publish(ctx, events.TemplateCreatedEvent{ID: t.ID, Code: t.Code})
}

func (p *TemplateEventPublisherKafka) PublishTemplateUpdated(ctx context.Context, t tplEntity.Template) error {
	return p.producer.Publish(ctx, events.TemplateCreatedEvent{ID: t.ID, Code: t.Code})
}

type VersionEventPublisherKafka struct {
	producer *libkafka.Producer[events.TemplateVersionCreatedEvent]
}

func NewVersionEventPublisherKafka(producer *libkafka.Producer[events.TemplateVersionCreatedEvent]) ucVer.VersionEventPublisher {
	return &VersionEventPublisherKafka{producer: producer}
}

func (p *VersionEventPublisherKafka) PublishVersionCreated(ctx context.Context, v verEntity.TemplateVersion) error {
	return p.producer.Publish(ctx, events.TemplateVersionCreatedEvent{
		ID: v.ID, TemplateID: v.TemplateID, Version: v.Version,
	})
}

func (p *VersionEventPublisherKafka) PublishVersionPublished(ctx context.Context, v verEntity.TemplateVersion) error {
	return p.producer.Publish(ctx, events.TemplateVersionCreatedEvent{
		ID: v.ID, TemplateID: v.TemplateID, Version: v.Version,
	})
}
