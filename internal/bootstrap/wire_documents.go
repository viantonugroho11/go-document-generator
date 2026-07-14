package bootstrap

import (
	"log"
	"strconv"
	"strings"

	documentsinfra "go-document-generator/internal/infrastructure/documents"
	kafkainfra "go-document-generator/internal/infrastructure/broker/kafka"
	miniostg "go-document-generator/internal/infrastructure/storage/minio"
	beginpg "go-document-generator/internal/repository/begin/postgres"
	cbpg "go-document-generator/internal/repository/documentcallbackattempts/postgres"
	logpg "go-document-generator/internal/repository/documentrenderlogs/postgres"
	tplpg "go-document-generator/internal/repository/documenttemplates/postgres"
	verpg "go-document-generator/internal/repository/documenttemplateversions/postgres"
	docpg "go-document-generator/internal/repository/documents/postgres"
	sharedStorage "go-document-generator/internal/shared/storage"
	"go-document-generator/internal/transport/apis"
	"go-document-generator/internal/transport/event/events"
	ucCb "go-document-generator/internal/usecase/documentcallbackattempts"
	ucDoc "go-document-generator/internal/usecase/documents"
	ucLog "go-document-generator/internal/usecase/documentrenderlogs"
	ucTpl "go-document-generator/internal/usecase/documenttemplates"
	ucVer "go-document-generator/internal/usecase/documenttemplateversions"

	cachetpl "go-document-generator/internal/infrastructure/cache/template"

	goredis "github.com/redis/go-redis/v9"
	libkafka "github.com/viantonugroho11/go-lib/kafka"
	"gorm.io/gorm"
)

func wireDocumentServices(db *gorm.DB, redis *goredis.Client) (apis.Services, func(), error) {
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

	rawTplRepo := tplpg.NewDocumentTemplatesRepository(db)
	rawVerRepo := verpg.NewDocumentTemplateVersionsRepository(db)

	tplRepo := cachetpl.NewCachedTemplateRepo(rawTplRepo, redis)
	verRepo := cachetpl.NewCachedVersionRepo(rawVerRepo, redis)
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
	docGeneratedProducer, err := libkafka.NewProducer[events.DocumentGeneratedEvent](
		c.KafkaBrokersList(), topicDoc,
		libkafka.WithKeyFunc(func(e events.DocumentGeneratedEvent) []byte { return []byte(e.RequestID) }),
	)
	if err != nil {
		_ = tplProducer.Close()
		_ = verProducer.Close()
		_ = docQueuedProducer.Close()
		_ = docRetryProducer.Close()
		return apis.Services{}, nil, err
	}
	docFailedProducer, err := libkafka.NewProducer[events.DocumentFailedEvent](
		c.KafkaBrokersList(), topicDoc,
		libkafka.WithKeyFunc(func(e events.DocumentFailedEvent) []byte { return []byte(e.RequestID) }),
	)
	if err != nil {
		_ = tplProducer.Close()
		_ = verProducer.Close()
		_ = docQueuedProducer.Close()
		_ = docRetryProducer.Close()
		_ = docGeneratedProducer.Close()
		return apis.Services{}, nil, err
	}
	docCancelledProducer, err := libkafka.NewProducer[events.DocumentCancelledEvent](
		c.KafkaBrokersList(), topicDoc,
		libkafka.WithKeyFunc(func(e events.DocumentCancelledEvent) []byte { return []byte(e.RequestID) }),
	)
	if err != nil {
		_ = tplProducer.Close()
		_ = verProducer.Close()
		_ = docQueuedProducer.Close()
		_ = docRetryProducer.Close()
		_ = docGeneratedProducer.Close()
		_ = docFailedProducer.Close()
		return apis.Services{}, nil, err
	}

	tplPublisher := kafkainfra.NewTemplateEventPublisherKafka(tplProducer)
	verPublisher := kafkainfra.NewVersionEventPublisherKafka(verProducer)
	docPublisher := kafkainfra.NewDocumentEventPublisherKafka(
		docQueuedProducer, docRetryProducer,
		docGeneratedProducer, docFailedProducer, docCancelledProducer,
	)
	selector := documentsinfra.NewSelector()

	// Storage provider: MinIO jika endpoint dikonfigurasi, fallback local.
	var storageProvider sharedStorage.Provider
	if c.Storage.Endpoint != "" && strings.ToLower(c.Storage.Provider) == "minio" {
		sp, spErr := miniostg.NewProvider(
			c.Storage.Endpoint, c.Storage.AccessKey, c.Storage.SecretKey,
			c.Storage.Bucket, c.Storage.UseSSL,
		)
		if spErr != nil {
			log.Printf("wire: minio init failed, fallback local: %v", spErr)
		} else {
			storageProvider = sp
		}
	}
	if storageProvider == nil {
		storageProvider = sharedStorage.NewLocalProvider(c.Storage.BaseDir)
	}

	svc := apis.Services{
		Templates:        ucTpl.NewService(tplRepo, tx, tplPublisher),
		TemplateVersions: ucVer.NewService(verRepo, tplRepo, tx, verPublisher),
		Documents:        ucDoc.NewService(docRepo, tplRepo, verRepo, tx, docPublisher, selector, storageProvider),
		RenderLogs:       ucLog.NewService(logRepo, docRepo),
		Callbacks:        ucCb.NewService(cbRepo, docRepo, c.Callback.HMACSecret),
	}

	cleanup := func() {
		_ = tplProducer.Close()
		_ = verProducer.Close()
		_ = docQueuedProducer.Close()
		_ = docRetryProducer.Close()
		_ = docGeneratedProducer.Close()
		_ = docFailedProducer.Close()
		_ = docCancelledProducer.Close()
	}
	return svc, cleanup, nil
}
