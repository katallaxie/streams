package writer

import (
	"context"
	"time"

	"github.com/katallaxie/pkg/logger"
	"github.com/segmentio/kafka-go"
)

// Writer is a Kafka writer.
type Writer interface {
	// Close closes the writer.
	Close() error

	// Stats returns the writer statistics.
	Stats() kafka.WriterStats

	// WriteMessage writes a message to Kafka.
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}

// Opt defines an option for a Kafka writer.
type Opt func(*kafka.Writer)

// WithTransport configures the transport for the writer.
func WithTransport(transport *kafka.Transport) Opt {
	return func(w *kafka.Writer) {
		w.Transport = transport
	}
}

// WithTopic configures the topic for the writer.
func WithTopic(topic string) Opt {
	return func(w *kafka.Writer) {
		w.Topic = topic
	}
}

// WithBrokers configures the brokers for the writer.
func WithBrokers(brokers ...string) Opt {
	return func(w *kafka.Writer) {
		w.Addr = kafka.TCP(brokers...)
	}
}

// WithBatchSize configures the batch size for the writer.
func WithBatchSize(batchSize int) Opt {
	return func(w *kafka.Writer) {
		w.BatchSize = batchSize
	}
}

// WithBatchTimeout configures the batch timeout for the writer.
func WithBatchTimeout(timeout time.Duration) Opt {
	return func(w *kafka.Writer) {
		w.BatchTimeout = timeout
	}
}

// WithBalancer configures the balancing mechanism for the writer.
func WithBalancer(b kafka.Balancer) Opt {
	return func(w *kafka.Writer) {
		w.Balancer = b
	}
}

// DefaultConfig returns the default configuration for a Kafka writer.
func DefaultConfig() *kafka.Writer {
	return &kafka.Writer{
		BatchSize:    1000,
		BatchTimeout: time.Second,
		Balancer:     &kafka.LeastBytes{},
		Logger:       kafka.LoggerFunc(logger.Infof),
		ErrorLogger:  kafka.LoggerFunc(logger.Errorf),
	}
}

// NewWriter creates a new Kafka writer.
func NewWriter(opts ...Opt) *kafka.Writer {
	w := DefaultConfig()

	for _, opt := range opts {
		opt(w)
	}

	return w
}
