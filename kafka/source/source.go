package source

import (
	"context"
	"sync"

	"github.com/ionos-cloud/streams/codec"
	"github.com/ionos-cloud/streams/msg"
	kgo "github.com/segmentio/kafka-go"
)

// Source is a Kafka source.
type Source[K, V any] struct {
	reader       *kgo.Reader
	ctx          context.Context
	keyDecoder   codec.Decoder[K]
	valueDecoder codec.Decoder[V]
	keyEncoder   codec.Encoder[K]

	err     error
	errOnce sync.Once
}

// WithContext is a constructor for a kafka source with a cancellation context.
func WithContext[K, V any](ctx context.Context, r *kgo.Reader, key codec.Decoder[K], value codec.Decoder[V], keyEncoder codec.Encoder[K]) *Source[K, V] {
	k := new(Source[K, V])
	k.ctx = ctx
	k.reader = r
	k.keyDecoder = key
	k.valueDecoder = value
	k.keyEncoder = keyEncoder

	return k
}

// Commit is a Kafka source commit.
func (k *Source[K, V]) Commit(msgs ...msg.Message[K, V]) error {
	mm := make([]kgo.Message, len(msgs))

	for i, m := range msgs {
		mm[i] = kgo.Message{Topic: k.reader.Config().Topic, Partition: m.Partition(), Offset: int64(m.Offset())}
	}

	return k.reader.CommitMessages(k.ctx, mm...)
}

// Message is a Kafka source message.
func (k *Source[K, V]) Messages() chan msg.Message[K, V] {
	out := make(chan msg.Message[K, V])

	go func() {
		for {
			m, err := k.reader.FetchMessage(k.ctx)
			if err != nil {
				k.fail(err)
				break
			}

			val, err := k.valueDecoder.Decode(m.Value)
			if err != nil {
				k.fail(err)
				break
			}

			key, err := k.keyDecoder.Decode(m.Key)
			if err != nil {
				k.fail(err)
				break
			}

			out <- msg.NewMessage(key, val, int(m.Offset), m.Partition, m.Topic)
		}

		close(out)
	}()

	return out
}

// Error is a Kafka source error.
func (k *Source[K, V]) Error() error {
	return k.err
}

func (k *Source[K, V]) fail(err error) {
	k.errOnce.Do(func() {
		k.err = err
	})
}
