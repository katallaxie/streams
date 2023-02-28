package source

import (
	"context"
	"errors"
	"io"
	"math"
	"sync"
	"time"

	"github.com/ionos-cloud/streams"
	"github.com/ionos-cloud/streams/codec"
	"github.com/ionos-cloud/streams/msg"

	kgo "github.com/segmentio/kafka-go"
)

var _ streams.Source[any, any] = (*Source[any, any])(nil)

// CommitMode is a Kafka source commit mode.
type CommitMode int

const (
	// CommitAuto commits messages automatically.
	CommitAuto CommitMode = iota

	// CommitManual commits messages manually.
	CommitManual
)

// Source is a Kafka source.
type Source[K, V any] struct {
	reader       *kgo.Reader
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
	commitMode    CommitMode
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

// DefaultOpts is the default options for a kafka source.
func DefaultOpts() *Opts {
	return &Opts{
		commitMode:    CommitManual,
		bufferSize:    1000,
		bufferTimeout: 1 * time.Second,
	}
}

// WithContext is a constructor for a kafka source with a cancellation context.
func WithContext[K, V any](ctx context.Context, r *kgo.Reader, key codec.Decoder[K], value codec.Decoder[V], keyEncoder codec.Encoder[K], opts ...Opt) *Source[K, V] {
	options := DefaultOpts()
	options.Configure(opts...)

	k := new(Source[K, V])

	k.ctx = ctx
	k.opts = options
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
	mark := make(msg.Marker[K, V])

	if s.opts.commitMode == CommitAuto {
		s.autoCommit(out)
	} else {
		s.manualCommit(out, mark)
		s.commit(mark)
	}

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

func (s *Source[K, V]) commit(mark msg.Marker[K, V]) {
	go func() {
		var buf []msg.Message[K, V]
		var count int

		ticker := time.NewTicker(math.MaxInt32 * time.Second)
		defer ticker.Stop()

		if s.opts.bufferTimeout > 0 {
			ticker.Reset(s.opts.bufferTimeout)
		} else {
			ticker.Stop()
		}

	LOOP:
		for {
			select {
			case <-ticker.C:
				err := s.Commit(buf...)
				if err != nil {
					s.fail(err)
					break LOOP
				}

				buf = buf[:0]
				count = 0
			case m, ok := <-mark:
				if !ok {
					break LOOP
				}

				buf = append(buf, m)
				count++

				if s.opts.bufferTimeout > 0 {
					continue
				}

				if count <= s.opts.bufferSize {
					continue
				}

				err := s.Commit(buf...)
				if err != nil {
					s.fail(err)
					break LOOP
				}

				buf = buf[:0]
				count = 0
			}
		}
	}()
}

func (s *Source[K, V]) manualCommit(out chan msg.Message[K, V], mark msg.Marker[K, V]) {
	go func() {
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

			out <- msg.NewMessage(key, val, int(m.Offset), m.Partition, m.Topic, mark)
		}

		close(out)
	}()
}

func (s *Source[K, V]) autoCommit(out chan msg.Message[K, V]) {
	go func() {
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

			out <- msg.NewMessage(key, val, int(m.Offset), m.Partition, m.Topic, nil)
		}

		close(out)
	}()
}
