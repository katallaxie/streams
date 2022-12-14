package reader

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

// Reader ...
type Reader interface {
	Close() error
	CommitMessages(ctx context.Context, msgs ...kafka.Message) error
	Config() kafka.ReaderConfig
	FetchMessage(ctx context.Context) (kafka.Message, error)
	Lag() int64
	Offset() int64
	ReadLag(ctx context.Context) (lag int64, err error)
	ReadMessage(ctx context.Context) (kafka.Message, error)
	SetOffset(offset int64) error
	SetOffsetAt(ctx context.Context, t time.Time) error
	Stats() kafka.ReaderStats
}

// Opt defines an option for a Kafka reader.
type Opt func(*kafka.ReaderConfig)

// WithDialer defines a dialer for the reader.
func WithDialer(dialer *kafka.Dialer) Opt {
	return func(r *kafka.ReaderConfig) {
		r.Dialer = dialer
	}
}

// WithGroupID defines the consumer group ID.
func WithGroupID(id string) Opt {
	return func(r *kafka.ReaderConfig) {
		r.GroupID = id
	}
}

// WithTopic defines the topic for the reader.
func WithTopic(topic string) Opt {
	return func(r *kafka.ReaderConfig) {
		r.Topic = topic
	}
}

// WithBrokers configures the brokers for the reader.
func WithBrokers(brokers ...string) Opt {
	return func(r *kafka.ReaderConfig) {
		r.Brokers = brokers
	}
}

// NewReader creates a new Kafka reader.
func NewReader(opts ...Opt) *kafka.Reader {
	cfg := kafka.ReaderConfig{
		CommitInterval: time.Second, // flushes commits to Kafka every second
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return kafka.NewReader(cfg)
}
