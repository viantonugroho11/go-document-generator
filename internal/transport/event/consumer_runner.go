package kafka

import (
	"context"

	"go-document-generator/internal/infrastructure/broker"
)

// RegisterConsumers menyiapkan dan menjalankan consumer Kafka dengan handler contoh.
// Anda bisa mengganti handler untuk memanggil usecase tertentu.
func RegisterConsumers(ctx context.Context, consumer broker.Consumer) {
	consumer.Start(ctx)
}
