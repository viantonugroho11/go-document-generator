package kafka

import (
	"context"
	"log"

	gokafka "github.com/viantonugroho11/go-lib/kafka"
)

// ExampleEvent adalah contoh event yang di-decode dari JSON message value.
type ExampleEvent struct {
	Message string `json:"message"`
}

// ExampleHandler memproses ExampleEvent.
// 1 consumer 1 handle: handler ini khusus untuk ExampleEvent saja.
type ExampleHandler struct{}

func (ExampleHandler) Name() string { return "example_consumer" }

func (ExampleHandler) Handle(ctx context.Context, evt ExampleEvent, _ ...gokafka.Header) gokafka.Progress {
	log.Printf("[example] message=%s", evt.Message)
	return gokafka.Progress{Status: gokafka.ProgressSuccess}
}

