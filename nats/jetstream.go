package nats

import (
	"context"
	"sync"

	"github.com/katallaxie/streams"
	"github.com/nats-io/nats.go"
)

var _ streams.Sourceable = (*JetStreamSource)(nil)

// JetStreamSourceConfig holds the configuration for a NATS JetStream source.
type JetStreamSourceConfig struct {
	Ack            bool
	AckOpts        []nats.AckOpt
	Conn           *nats.Conn
	ConsumerName   string
	FetchBatchSize int
	JetStreamCtx   nats.JetStreamContext
	PullOpts       []nats.PullOpt
	Subject        string
	SubOpts        []nats.SubOpt
}

// DefaultJetStreamSourceConfig returns a default JetStream source configuration.
func DefaultJetStreamSourceConfig() *JetStreamSourceConfig {
	return &JetStreamSourceConfig{
		FetchBatchSize: 10,
		Ack:            true,
		SubOpts:        []nats.SubOpt{},
		PullOpts:       []nats.PullOpt{},
		AckOpts:        []nats.AckOpt{},
	}
}

// JetStreamSource represents a NATS JetStream.
type JetStreamSource struct {
	out     chan any
	sub     *nats.Subscription
	err     error
	errOnce sync.Once
	cfg     *JetStreamSourceConfig
}

// NewJetStreamSource returns a new JetStreamSource connector that reads messages from a NATS JetStream subscription.
func NewJetStreamSource(ctx context.Context, cfg *JetStreamSourceConfig) (*JetStreamSource, error) {
	jetstreamSource := &JetStreamSource{
		out: make(chan any),
		cfg: cfg,
	}

	go jetstreamSource.attach(ctx)

	return jetstreamSource, nil
}

// Error returns the error.
func (j *JetStreamSource) Error() error {
	return j.err
}

func (j *JetStreamSource) fail(err error) {
	j.errOnce.Do(func() {
		j.err = err
	})
}

// Pipe pipes the output channel of the ReaderSource connector to the input channel.
func (j *JetStreamSource) Pipe(operator streams.Operatable) streams.Operatable {
	streams.Pipe(j, operator)
	return operator
}

// Out returns the output channel of the ReaderSource connector.
func (j *JetStreamSource) Out() <-chan any {
	return j.out
}

func (j *JetStreamSource) attach(ctx context.Context) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		default:
		}
		messages, err := j.sub.Fetch(j.cfg.FetchBatchSize, j.cfg.PullOpts...)
		if err != nil {
			j.fail(err)
			break loop
		}

		if len(messages) == 0 {
			continue
		}

		for _, msg := range messages {
			j.out <- msg
			// TODO: Handle the message, e.g., process it or transform it
			// Acknowledge the message if Ack is enabled
		}
	}

	if err := j.sub.Drain(); err != nil {
		j.fail(err)
	}
	close(j.out)
}
