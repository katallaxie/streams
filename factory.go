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
type StreamImpl struct {
	in    chan msg.Message
	mark  chan msg.Message
	close chan bool
	err   chan error
	opts  *Opts
}

// NewStream from a source of messages.
func NewStream(src Source, opts ...Opt) *StreamImpl {
	options := new(Opts)
	options.Configure(opts...)

	stream := new(StreamImpl)
	stream.opts = options
	stream.mark = make(chan msg.Message)
	stream.in = src.Messages()

	go func() {
		var count int
		var buf []msg.Message

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
