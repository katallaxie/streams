package sink

import (
	"context"

	"github.com/katallaxie/streams"
	"github.com/katallaxie/streams/codec"
	"github.com/katallaxie/streams/msg"
	"github.com/nats-io/nats.go"
)

var _ streams.Sink[any, any] = (*Sink[any, any])(nil)

// Sink is a nats sink.
type Sink[K, V any] struct {
	conn         *nats.Conn
	ctx          context.Context
	keyEncoder   codec.Encoder[K]
	valueEncoder codec.Encoder[V]
}

// WithContext is a constructor for a kafka sink with a cancellation context.
func WithContext[K, V any](ctx context.Context, conn *nats.Conn, key codec.Encoder[K], value codec.Encoder[V]) *Sink[K, V] {
	k := new(Sink[K, V])
	k.ctx = ctx
	k.conn = conn
	k.keyEncoder = key
	k.valueEncoder = value

	return k
}

// NewSink is a nats sink constructor.
func NewSink[K, V any]() *Sink[K, V] {
	n := new(Sink[K, V])

	return n
}

// Write is a nats sink writer.
func (s *Sink[K, V]) Write(messages ...msg.Message[K, V]) error {
	for _, m := range messages {
		val, err := s.valueEncoder.Encode(m.Value())
		if err != nil {
			return err
		}

		err = s.conn.Publish(m.Topic(), val)
		if err != nil {
			return err
		}
	}

	return nil
}
