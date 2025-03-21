package source

import (
	"context"
	"sync"

	"github.com/katallaxie/pkg/cast"
	"github.com/katallaxie/streams"
	"github.com/katallaxie/streams/codec"
	"github.com/katallaxie/streams/msg"

	"github.com/nats-io/nats.go"
)

var _ streams.Source[any, any] = (*Source[any, any])(nil)

// Source is a source of NATS messages.
type Source[K, V any] struct {
	ctx          context.Context
	keyDecoder   codec.Decoder[K]
	valueDecoder codec.Decoder[V]

	sub *nats.Subscription

	opts *Opts

	err     error
	errOnce sync.Once
}

// Opts is a set of options for a NATS source.
type Opts struct{}

// Opt is a NATS source option.
type Opt func(o *Opts)

// Configure is a function that configures a NATS source.
func (o *Opts) Configure(opts ...Opt) {
	for _, opt := range opts {
		opt(o)
	}
}

// DefaultOpts is a function that returns a set of default options for a NATS source.
func DefaultOpts() *Opts {
	return &Opts{}
}

// WithContext is a function that configures a NATS source with a context.
func WithContext[K, V any](ctx context.Context, sub *nats.Subscription, key codec.Decoder[K], value codec.Decoder[V], opts ...Opt) *Source[K, V] {
	options := DefaultOpts()
	options.Configure(opts...)

	k := new(Source[K, V])
	k.ctx = ctx
	k.keyDecoder = key
	k.valueDecoder = value
	k.sub = sub
	k.opts = options

	return k
}

// Messages is a function that returns a channel of NATS messages.
func (s *Source[K, V]) Messages() chan msg.Message[K, V] {
	out := make(chan msg.Message[K, V])

	go func(sub *nats.Subscription) {
		for {

			m, err := sub.NextMsgWithContext(s.ctx)
			if err != nil {
				s.fail(err)
				break
			}

			val, err := s.valueDecoder.Decode(m.Data)
			if err != nil {
				s.fail(err)
				break
			}

			out <- msg.NewMessage(cast.Zero[K](), val, 0, 0, m.Subject, nil)
		}

		close(out)
	}(s.sub)

	return out
}

// Commit is a function that commits a NATS message.
func (s *Source[K, V]) Commit(_ ...msg.Message[K, V]) error {
	return nil
}

// Error is a function that returns the error of a NATS source.
func (s *Source[K, V]) Error() error {
	return s.err
}

func (s *Source[K, V]) fail(err error) {
	s.errOnce.Do(func() {
		s.err = err
	})
}
