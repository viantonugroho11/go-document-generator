package bootstrap

import (
	"log"
	"strconv"
	"strings"

	documentsinfra "go-document-generator/internal/infrastructure/documents"
	kafkainfra "go-document-generator/internal/infrastructure/broker/kafka"
	miniostg "go-document-generator/internal/infrastructure/storage/minio"
	ossstg  "go-document-generator/internal/infrastructure/storage/oss"
	s3stg   "go-document-generator/internal/infrastructure/storage/s3"
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
	topicDocProcess := c.Kafka.TopicDocumentProcess
	if topicDocProcess == "" {
		topicDocProcess = "document-process"
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

	// document-events: lifecycle events (DocumentEvent) — partition by request_id.
	docEventProducer, err := libkafka.NewProducer[events.DocumentEvent](
		c.KafkaBrokersList(), topicDoc,
		libkafka.WithKeyFunc(func(e events.DocumentEvent) []byte { return []byte(e.ResourceID) }),
	)
	if err != nil {
		_ = tplProducer.Close()
		_ = verProducer.Close()
		return apis.Services{}, nil, err
	}

	// document-events: bulk ops (DocumentBulkEvent) — partition by output_path.
	docBulkProducer, err := libkafka.NewProducer[events.DocumentBulkEvent](
		c.KafkaBrokersList(), topicDoc,
		libkafka.WithKeyFunc(func(e events.DocumentBulkEvent) []byte { return []byte(e.ResourceID) }),
	)
	if err != nil {
		_ = tplProducer.Close()
		_ = verProducer.Close()
		_ = docEventProducer.Close()
		return apis.Services{}, nil, err
	}

	// document-process: generation trigger — partition by request_id.
	docProcessProducer, err := libkafka.NewProducer[events.DocumentEvent](
		c.KafkaBrokersList(), topicDocProcess,
		libkafka.WithKeyFunc(func(e events.DocumentEvent) []byte { return []byte(e.ResourceID) }),
	)
	if err != nil {
		_ = tplProducer.Close()
		_ = verProducer.Close()
		_ = docEventProducer.Close()
		_ = docBulkProducer.Close()
		return apis.Services{}, nil, err
	}

	tplPublisher := kafkainfra.NewTemplateEventPublisherKafka(tplProducer)
	verPublisher := kafkainfra.NewVersionEventPublisherKafka(verProducer)
	docPublisher := kafkainfra.NewDocumentEventPublisherKafka(docEventProducer, docBulkProducer, docProcessProducer)
	selector := documentsinfra.NewSelector()

	// Storage provider dipilih berdasarkan config storage.provider.
	var storageProvider sharedStorage.Provider
	if c.Storage.Endpoint != "" {
		var sp sharedStorage.Provider
		var spErr error
		switch strings.ToLower(c.Storage.Provider) {
		case "s3", "aws":
			sp, spErr = s3stg.NewProvider(
				c.Storage.Endpoint, c.Storage.AccessKey, c.Storage.SecretKey,
				c.Storage.Bucket, c.Storage.UseSSL,
			)
		case "oss", "alibaba":
			sp, spErr = ossstg.NewProvider(
				c.Storage.Endpoint, c.Storage.AccessKey, c.Storage.SecretKey,
				c.Storage.Bucket, c.Storage.UseSSL,
			)
		case "minio":
			sp, spErr = miniostg.NewProvider(
				c.Storage.Endpoint, c.Storage.AccessKey, c.Storage.SecretKey,
				c.Storage.Bucket, c.Storage.UseSSL,
			)
		default:
			log.Printf("wire: storage provider %q tidak dikenal, fallback local", c.Storage.Provider)
		}
		if spErr != nil {
			log.Printf("wire: storage %s init failed, fallback local: %v", c.Storage.Provider, spErr)
		} else if sp != nil {
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
		_ = docEventProducer.Close()
		_ = docBulkProducer.Close()
		_ = docProcessProducer.Close()
	}
	return svc, cleanup, nil
}
