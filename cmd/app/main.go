package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"go-document-generator/internal/config"
	kafkainfra "go-document-generator/internal/infrastructure/broker/kafka"
	redisinfra "go-document-generator/internal/infrastructure/cache/redis"
	pginfra "go-document-generator/internal/infrastructure/database/postgres"
	tmplpg "go-document-generator/internal/repository/documenttemplates/postgres"
	verpg "go-document-generator/internal/repository/documenttemplateversions/postgres"
	userpg "go-document-generator/internal/repository/user/postgres"
	"go-document-generator/internal/transport/apis"
	kafkarunner "go-document-generator/internal/transport/event/kafka"
	usecasetmpl "go-document-generator/internal/usecase/documenttemplates"
	usecasever "go-document-generator/internal/usecase/documenttemplateversions"
	usecaseusers "go-document-generator/internal/usecase/users"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	confLoader "github.com/viantonugroho11/go-lib/config"
)

func main() {
	cfg := config.Configuration{}
	loader := confLoader.New(
		"",                      // ENV prefix
		"go-document-generator", // Consul KV key
		os.Getenv("CONSUL_URL"),
		confLoader.WithConfigFileSearchPaths("./config"),
	)
	err := loader.Load(&cfg)
	if err != nil {
		os.Exit(1)
	}

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// Wiring dependencies
	ctx := context.Background()
	db, err := pginfra.Connect(ctx, cfg.PGDSN())
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err := pginfra.Migrate(db); err != nil {
		log.Fatalf("db migrate error: %v", err)
	}
	userRepo := userpg.NewUserRepository(db)
	userService := usecaseusers.NewUserService(userRepo)
	tmplRepo := tmplpg.NewDocumentTemplatesRepository(db)
	verRepo := verpg.NewDocumentTemplateVersionsRepository(db)
	tmplService := usecasetmpl.NewTemplatesService(tmplRepo, verRepo)
	verService := usecasever.NewVersionsService(verRepo)
	apis.RegisterRoutes(e, userService, tmplService, verService)

	// Init Redis
	redisClient, err := redisinfra.NewClient(cfg.Redis.Addr, cfg.Redis.Password, strconv.Itoa(cfg.Redis.DB))
	if err != nil {
		log.Fatalf("redis init error: %v", err)
	}
	defer redisClient.Close()

	// Init Kafka Consumer via go-lib adapter
	consumer, err := kafkainfra.NewGoLibConsumer[kafkarunner.ExampleEvent](
		cfg.KafkaBrokersList(),
		cfg.Kafka.GroupID,
		cfg.Kafka.Topic,
		kafkarunner.ExampleHandler{},
	)
	if err != nil {
		log.Fatalf("kafka consumer init error: %v", err)
	}
	kafkarunner.RegisterConsumers(ctx, consumer)
	defer consumer.Close()

	// HTTP server with graceful shutdown
	server := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      e,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := e.StartServer(server); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()
	log.Printf("server listening on :%s", cfg.App.Port)

	// wait for interrupt signal to gracefully shutdown the server with a timeout
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	} else {
		log.Println("server shutdown gracefully")
	}
}
