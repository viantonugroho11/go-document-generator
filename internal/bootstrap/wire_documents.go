package bootstrap

import (
	"strconv"

	documentsinfra "go-document-generator/internal/infrastructure/documents"
	kafkainfra "go-document-generator/internal/infrastructure/broker/kafka"
	beginpg "go-document-generator/internal/repository/begin/postgres"
	cbpg "go-document-generator/internal/repository/documentcallbackattempts/postgres"
	logpg "go-document-generator/internal/repository/documentrenderlogs/postgres"
	tplpg "go-document-generator/internal/repository/documenttemplates/postgres"
	verpg "go-document-generator/internal/repository/documenttemplateversions/postgres"
	docpg "go-document-generator/internal/repository/documents/postgres"
	"go-document-generator/internal/transport/apis"
	"go-document-generator/internal/transport/event/events"
	ucCb "go-document-generator/internal/usecase/documentcallbackattempts"
	ucDoc "go-document-generator/internal/usecase/documents"
	ucLog "go-document-generator/internal/usecase/documentrenderlogs"
	ucTpl "go-document-generator/internal/usecase/documenttemplates"
	ucVer "go-document-generator/internal/usecase/documenttemplateversions"

	libkafka "github.com/viantonugroho11/go-lib/kafka"
	"gorm.io/gorm"
)

func wireDocumentServices(db *gorm.DB) (apis.Services, func(), error) {
	c := Config()
	topicTpl := c.Kafka.TopicTemplates
	if topicTpl == "" {
		topicTpl = "template-events"
	}
	topicVer := c.Kafka.TopicTemplateVersions
	if topicVer == "" {
		topicVer = "template-version-events"
	}
	topicDoc := c.Kafka.TopicDocuments
	if topicDoc == "" {
		topicDoc = "document-events"
	}

	tx := beginpg.NewBeginRepository(db)

	tplRepo := tplpg.NewDocumentTemplatesRepository(db)
	verRepo := verpg.NewDocumentTemplateVersionsRepository(db)
	docRepo := docpg.NewDocumentsRepository(db)
	logRepo := logpg.NewDocumentRenderLogsRepository(db)
	cbRepo := cbpg.NewDocumentCallbackAttemptsRepository(db)

	tplProducer, err := libkafka.NewProducer[events.TemplateCreatedEvent](
		c.KafkaBrokersList(), topicTpl,
		libkafka.WithKeyFunc(func(e events.TemplateCreatedEvent) []byte { return []byte(e.Code) }),
	)
	if err != nil {
		return apis.Services{}, nil, err
	}
	verProducer, err := libkafka.NewProducer[events.TemplateVersionCreatedEvent](
		c.KafkaBrokersList(), topicVer,
		libkafka.WithKeyFunc(func(e events.TemplateVersionCreatedEvent) []byte {
			return []byte(strconv.FormatInt(e.TemplateID, 10))
		}),
	)
	if err != nil {
		_ = tplProducer.Close()
		return apis.Services{}, nil, err
	}
	docQueuedProducer, err := libkafka.NewProducer[events.DocumentQueuedEvent](
		c.KafkaBrokersList(), topicDoc,
		libkafka.WithKeyFunc(func(e events.DocumentQueuedEvent) []byte { return []byte(e.RequestID) }),
	)
	if err != nil {
		_ = tplProducer.Close()
		_ = verProducer.Close()
		return apis.Services{}, nil, err
	}
	docRetryProducer, err := libkafka.NewProducer[events.DocumentRetriedEvent](
		c.KafkaBrokersList(), topicDoc,
		libkafka.WithKeyFunc(func(e events.DocumentRetriedEvent) []byte { return []byte(e.RequestID) }),
	)
	if err != nil {
		_ = tplProducer.Close()
		_ = verProducer.Close()
		_ = docQueuedProducer.Close()
		return apis.Services{}, nil, err
	}

	tplPublisher := kafkainfra.NewTemplateEventPublisherKafka(tplProducer)
	verPublisher := kafkainfra.NewVersionEventPublisherKafka(verProducer)
	docPublisher := kafkainfra.NewDocumentEventPublisherKafka(docQueuedProducer, docRetryProducer)
	selector := documentsinfra.NewSelector()

	svc := apis.Services{
		Templates:        ucTpl.NewService(tplRepo, tx, tplPublisher),
		TemplateVersions: ucVer.NewService(verRepo, tplRepo, tx, verPublisher),
		Documents:        ucDoc.NewService(docRepo, tplRepo, verRepo, tx, docPublisher, selector),
		RenderLogs:       ucLog.NewService(logRepo, docRepo),
		Callbacks:        ucCb.NewService(cbRepo, docRepo),
	}

	cleanup := func() {
		_ = tplProducer.Close()
		_ = verProducer.Close()
		_ = docQueuedProducer.Close()
		_ = docRetryProducer.Close()
	}
	return svc, cleanup, nil
}
