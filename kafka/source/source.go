package source

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/ionos-cloud/streams"
	"github.com/ionos-cloud/streams/codec"
	"github.com/ionos-cloud/streams/msg"
	kgo "github.com/segmentio/kafka-go"
)

// CommitMode is a commit mode.
type CommitMode int

const (
	// ManualCommit is a manual commit.
	ManualCommit CommitMode = iota

	// AutoCommit is an auto commit.
	AutoCommit
)

// Source is a Kafka source.
type Source[K, V any] struct {
	reader       *kgo.Reader
	out          streams.MessageReceiver[K, V]
	ctx          context.Context
	keyDecoder   codec.Decoder[K]
	valueDecoder codec.Decoder[V]
	keyEncoder   codec.Encoder[K]

	opts *Opts

	err     error
	errOnce sync.Once
}

// Opts is a set of options for a kafka source.
type Opts struct {
	mode          CommitMode
	bufferSize    int
	bufferTimeout time.Duration
}

// Opt is a Kafka source option.
type Opt func(o *Opts)

// Configure is a function that configures a kafka source.
func (o *Opts) Configure(opts ...Opt) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithBufferSize configures the buffer size for a kafka source.
func WithBufferSize(size int) Opt {
	return func(o *Opts) {
		o.bufferSize = size
	}
}

// WithAutoCommit configures the kafka source to auto commit.
func WithAutoCommit() Opt {
	return func(o *Opts) {
		o.mode = AutoCommit
	}
}

// WithManualCommit configures the kafka source to manual commit.
func WithManualCommit() Opt {
	return func(o *Opts) {
		o.mode = ManualCommit
	}
}

// DefaultOpts is the default options for a kafka source.
func DefaultOpts() *Opts {
	return &Opts{
		mode:          AutoCommit,
		bufferSize:    100,
		bufferTimeout: 1 * time.Second,
	}
}

// WithContext is a constructor for a kafka source with a cancellation context.
func WithContext[K, V any](ctx context.Context, r *kgo.Reader, key codec.Decoder[K], value codec.Decoder[V], keyEncoder codec.Encoder[K], opts ...Opt) *Source[K, V] {
	options := DefaultOpts()
	options.Configure(opts...)

	k := new(Source[K, V])

	k.out = make(streams.MessageReceiver[K, V])
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

	if s.opts.mode == AutoCommit {
		go s.autoCommit(out)()
	} else {
		go s.manualCommit(out)()
	}

	return out
}

// Close is a Kafka source close.
func (s *Source[K, V]) Close() error {
	return s.reader.Close()
}

// Error is a Kafka source error.
func (s *Source[K, V]) Error() error {
	return s.err
}

func (s *Source[K, V]) manualCommit(out chan msg.Message[K, V]) func() {
	return func() {
		for {
			m, err := s.reader.FetchMessage(s.ctx)
			if errors.Is(err, io.EOF) {
				break
			}

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

			out <- msg.NewMessage(key, val, int(m.Offset), m.Partition, m.Topic)
		}

		close(out)
	}
}

func (s *Source[K, V]) autoCommit(out chan msg.Message[K, V]) func() {
	return func() {
		for {
			m, err := s.reader.ReadMessage(s.ctx)
			if errors.Is(err, io.EOF) {
				break
			}

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

			out <- msg.NewMessage(key, val, int(m.Offset), m.Partition, m.Topic)
		}

		close(out)
	}
}

func (s *Source[K, V]) fail(err error) {
	s.errOnce.Do(func() {
		s.err = err
	})
}
