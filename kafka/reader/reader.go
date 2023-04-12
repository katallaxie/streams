package reader

import (
	"context"
	"time"

	"github.com/katallaxie/pkg/logger"
	"github.com/segmentio/kafka-go"
)

// Reader ...
type Reader interface {
	// Close closes the reader.
	Close() error

	// CommitMessages commits the messages to Kafka.
	CommitMessages(ctx context.Context, msgs ...kafka.Message) error

	// Config returns the reader configuration.
	Config() kafka.ReaderConfig

	// FetchMessage fetches a message from Kafka.
	FetchMessage(ctx context.Context) (kafka.Message, error)

	// Lag returns the current lag of the reader.
	Lag() int64

	// Offset returns the current offset of the reader.
	Offset() int64

	// ReadLag reads the current lag of the reader.
	ReadLag(ctx context.Context) (lag int64, err error)

	// ReadMessage reads a message from Kafka.
	ReadMessage(ctx context.Context) (kafka.Message, error)

	// SetOffset sets the offset of the reader.
	SetOffset(offset int64) error

	// SetOffsetAt sets the offset of the reader at the given time.
	SetOffsetAt(ctx context.Context, t time.Time) error

	// Stats returns the reader statistics.
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

// WithLogger configures the logger for the reader.
func WithLogger(logger kafka.Logger) Opt {
	return func(rc *kafka.ReaderConfig) {
		rc.Logger = logger
	}
}

// WithErrorLogger configures the error logger for the reader.
func WithErrorLogger(logger kafka.Logger) Opt {
	return func(rc *kafka.ReaderConfig) {
		rc.ErrorLogger = logger
	}
}

// DefaultConfig returns the default configuration for a Kafka reader.
func DefaultConfig() kafka.ReaderConfig {
	return kafka.ReaderConfig{
		Logger:         kafka.LoggerFunc(logger.Infof),
		ErrorLogger:    kafka.LoggerFunc(logger.Errorf),
		CommitInterval: 0,
	}
}

// NewReader creates a new Kafka reader.
func NewReader(opts ...Opt) *kafka.Reader {
	cfg := DefaultConfig()

	for _, opt := range opts {
		opt(&cfg)
	}

	return kafka.NewReader(cfg)
}
