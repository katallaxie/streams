package sink

import (
	"context"

	"github.com/ionos-cloud/streams"
	"github.com/ionos-cloud/streams/codec"
	"github.com/ionos-cloud/streams/msg"

	kgo "github.com/segmentio/kafka-go"
)

var _ streams.Sink[any, any] = (*Sink[any, any])(nil)

// Sink is a noop sink.
type Sink[K, V any] struct {
	writer       *kgo.Writer
	ctx          context.Context
	keyEncoder   codec.Encoder[K]
	valueEncoder codec.Encoder[V]
}

// WithContext is a constructor for a kafka sink with a cancellation context.
func WithContext[K, V any](ctx context.Context, w *kgo.Writer, key codec.Encoder[K], value codec.Encoder[V]) *Sink[K, V] {
	k := new(Sink[K, V])
	k.ctx = ctx
	k.writer = w
	k.keyEncoder = key
	k.valueEncoder = value

	return k
}

// NewSink is a noop sink constructor.
func NewSink[K, V any]() *Sink[K, V] {
	n := new(Sink[K, V])

	return n
}

// Write is a noop sink writer.
func (s *Sink[K, V]) Write(messages ...msg.Message[K, V]) error {
	mm := make([]kgo.Message, len(messages))

	for i, m := range messages {
		val, err := s.valueEncoder.Encode(m.Value())
		if err != nil {
			return err
		}

		key, err := s.keyEncoder.Encode(m.Key())
		if err != nil {
			return err
		}

		mm[i] = kgo.Message{
			Key:   key,
			Value: val,
			Topic: m.Topic(),
		}
	}

	return s.writer.WriteMessages(s.ctx, mm...)
}
