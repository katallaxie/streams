package kafka

import (
	"context"

	"github.com/ionos-cloud/streams/msg"
	kgo "github.com/segmentio/kafka-go"
)

type kafka[T any] struct {
	reader       *kgo.Reader
	ctx          context.Context
	keyDecoder   Decoder[T]
	valueDecoder Decoder[T]
}

// WithContext is a constructor for a kafka source with a cancellation context.
func WithContext[T any](ctx context.Context, r *kgo.Reader) *kafka[T] {
	k := new(kafka[T])
	k.ctx = ctx
	k.reader = r

	return k
}

// Commit ...
func (k *kafka[T]) Commit(msgs ...msg.Message) error {
	mm := make([]kgo.Message, len(msgs))
	for i, m := range msgs {
		mm[i] = kgo.Message{Key: []byte(m.Key())}
	}

	return k.reader.CommitMessages(k.ctx, mm...)
}

// Message ...
func (k *kafka[T]) Messages() chan msg.Message {
	out := make(chan msg.Message)

	go func() {
		for {
			m, err := k.reader.FetchMessage(k.ctx)
			if err != nil {
				break
			}

			out <- msg.NewMessage(string(m.Key))
		}

		close(out)
	}()

	return out
}
