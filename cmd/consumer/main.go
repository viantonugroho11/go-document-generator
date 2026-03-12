package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-document-generator/internal/config"
	"go-document-generator/internal/infrastructure/broker"
	kafkainfra "go-document-generator/internal/infrastructure/broker/kafka"
	kafkarunner "go-document-generator/internal/transport/event/kafka"
	consumerrunner "go-document-generator/internal/transport/event"

	confLoader "github.com/viantonugroho11/go-lib/config"
)

func main() {
	// Pilih consumer via flag: -consumer=example
	var consumerName string
	flag.StringVar(&consumerName, "consumer", "example", "nama consumer yang akan dijalankan")
	flag.Parse()

	// Load configuration (Consul/env/file)
	cfg := config.Configuration{}
	loader := confLoader.New(
		"",                      // ENV prefix (kosong => disable)
		"go-document-generator", // Consul KV key (opsional)
		os.Getenv("CONSUL_URL"), // Consul address (opsional)
		confLoader.WithConfigFileSearchPaths("./config"),
	)
	if err := loader.Load(&cfg); err != nil {
		log.Fatalf("config load error: %v", err)
	}

	// Root context with cancellation on signal
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize & start Kafka Consumer via go-lib adapter (1 consumer 1 handler)
	var (
		consumer broker.Consumer
		err      error
	)
	switch consumerName {
	case "example":
		consumer, err = kafkainfra.NewGoLibConsumer[kafkarunner.ExampleEvent](
			cfg.KafkaBrokersList(),
			cfg.Kafka.GroupID,
			cfg.Kafka.Topic,
			kafkarunner.ExampleHandler{},
		)
		if err != nil {
			log.Fatalf("kafka consumer init error: %v", err)
		}
	default:
		log.Fatalf("unknown consumer: %s", consumerName)
	}

	// Ensure close outside switch and single registration/logging
	defer func() {
		if consumer != nil {
			if cerr := consumer.Close(); cerr != nil {
				log.Printf("kafka consumer close error: %v", cerr)
			}
		}
	}()
	consumerrunner.RegisterConsumers(ctx, consumer)
	log.Printf("consumer started: %s group=%s topic=%s", consumerName, cfg.Kafka.GroupID, cfg.Kafka.Topic)

	// Wait for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("shutdown signal received, stopping consumer...")
}
