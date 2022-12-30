package streams

import (
	"github.com/ionos-cloud/streams/msg"
)

// Opts is a set of options for a stream.
type Opts struct {
	buffer int
}

// Comfigure is a function that configures a stream.
func (o *Opts) Configure(opts ...Opt) {
	for _, opt := range opts {
		opt(o)
	}
}

// Opt is a function that configures a stream.
type Opt func(*Opts)

// WithBuffer configures the buffer size for a stream.
func WithBuffer(size int) Opt {
	return func(o *Opts) {
		o.buffer = size
	}
}

// StreamImpl implements Stream.
type StreamImpl[K, V any] struct {
	in    chan msg.Message[K, V]
	mark  chan msg.Message[K, V]
	close chan bool
	err   chan error
	opts  *Opts
}

// NewStream from a source of messages.
func NewStream[K, V any](src Source[K, V], opts ...Opt) *StreamImpl[K, V] {
	options := new(Opts)
	options.Configure(opts...)

	stream := new(StreamImpl[K, V])
	stream.opts = options
	stream.mark = make(chan msg.Message[K, V])
	stream.in = src.Messages()

	go func() {
		var count int
		var buf []msg.Message[K, V]

		for m := range stream.mark {
			if m.Marked() {
				continue
			}

			buf = append(buf, m)
			count++

			m.Mark()

			if count <= stream.opts.buffer {
				continue
			}

			err := src.Commit(buf...)
			if err != nil {
				stream.Fail(err)
				return
			}

			buf = buf[:0]
			count = 0
		}
	}()

	return stream
}
