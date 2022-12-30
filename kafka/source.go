package kafka

import (
	"context"

	"github.com/ionos-cloud/streams/msg"
	kgo "github.com/segmentio/kafka-go"
)

type kafka[K, V any] struct {
	reader       *kgo.Reader
	ctx          context.Context
	keyDecoder   Decoder[K]
	valueDecoder Decoder[V]
}

// WithContext is a constructor for a kafka source with a cancellation context.
func WithContext[K, V any](ctx context.Context, r *kgo.Reader) *kafka[K, V] {
	k := new(kafka[K, V])
	k.ctx = ctx
	k.reader = r

	return k
}

// Commit ...
func (k *kafka[K, V]) Commit(msgs ...msg.Message[K, V]) error {
	mm := make([]kgo.Message, len(msgs))
	for i := range msgs {
		mm[i] = kgo.Message{Key: []byte("help")}
	}

	return k.reader.CommitMessages(k.ctx, mm...)
}

// Message ...
func (k *kafka[K, V]) Messages() chan msg.Message[K, V] {
	out := make(chan msg.Message[K, V])

	go func() {
		for {
			m, err := k.reader.FetchMessage(k.ctx)
			if err != nil {
				break
			}

			val, err := k.valueDecoder.Decode(m.Value)
			if err != nil {
				break
			}

			key, err := k.keyDecoder.Decode(m.Key)
			if err != nil {
				break
			}

			out <- msg.NewMessage(key, val)
		}

		close(out)
	}()

	return out
}
