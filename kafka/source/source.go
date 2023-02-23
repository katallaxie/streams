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

	mark msg.Marker[K, V]

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
func (s *Source[K, V]) Commit(msgs ...msg.Message[K, V]) error {
	mm := make([]kgo.Message, len(msgs))

	for i, m := range msgs {
		mm[i] = kgo.Message{Topic: s.reader.Config().Topic, Partition: m.Partition(), Offset: int64(m.Offset())}
	}

	return s.reader.CommitMessages(s.ctx, mm...)
}

// Message is a Kafka source message.
func (s *Source[K, V]) Messages() chan msg.Message[K, V] {
	out := make(chan msg.Message[K, V])

	go func() {
		for {
			m, err := s.reader.FetchMessage(s.ctx)
			if err != nil {
				s.fail(err)
				break
			}

			val, err := s.valueDecoder.Decode(m.Value)
			if err != nil {
				s.fail(err)
				break
			}

			key, err := s.keyDecoder.Decode(m.Key)
			if err != nil {
				s.fail(err)
				break
			}

			out <- msg.NewMessage(key, val, int(m.Offset), m.Partition, m.Topic, s.mark)
		}

		close(out)
	}()

	return out
}

// Error is a Kafka source error.
func (s *Source[K, V]) Error() error {
	return s.err
}

func (s *Source[K, V]) fail(err error) {
	s.errOnce.Do(func() {
		s.err = err
	})
}
