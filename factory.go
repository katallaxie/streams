package streams

import (
	"context"

	"github.com/ionos-cloud/streams/msg"
)

// Option ...
type Option interface{}

// Iterable ...
type Iterable interface {
	Observe(opts ...Option) <-chan *msg.Message
}

// StreamImpl implements Stream.
type StreamImpl struct {
	parent   context.Context
	iterable Iterable
}

// Observable ...
type Observable interface{}

// FromChannel ...
func FromChannel(ctx context.Context, next <-chan *msg.Message, opts ...Option) Observable {
	return &StreamImpl{
		parent:   ctx,
		iterable: newChannelIterable(next, opts...),
	}
}
