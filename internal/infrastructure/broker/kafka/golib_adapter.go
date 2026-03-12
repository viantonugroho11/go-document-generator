package kafka

import (
	"context"

	gokafka "github.com/viantonugroho11/go-lib/kafka"
)

// GoLibConsumer adalah adapter tipis agar consumer dari go-lib/kafka
// bisa mengikuti kontrak broker.Consumer kita (Start/Close).
// Gunakan tipe event E sesuai payload JSON Anda.
type GoLibConsumer[E any] struct {
	inner gokafka.Consumer
}

// NewGoLibConsumer membuat consumer golib yang di-adapt ke interface broker.Consumer.
func NewGoLibConsumer[E any](
	brokers []string,
	groupID string,
	topic string,
	handler gokafka.EventHandler[E],
	opts ...gokafka.ConsumerOption,
) (*GoLibConsumer[E], error) {
	c, err := gokafka.NewConsumer[E](brokers, groupID, topic, handler, opts...)
	if err != nil {
		return nil, err
	}
	return &GoLibConsumer[E]{inner: c}, nil
}

// Start menjalankan konsumsi pesan.
func (c *GoLibConsumer[E]) Start(ctx context.Context) {
	c.inner.Start(ctx)
}

// Close menutup consumer dan resource terkait.
func (c *GoLibConsumer[E]) Close() error {
	return c.inner.Close()
}

