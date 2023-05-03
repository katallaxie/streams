package noop

import (
	"github.com/katallaxie/streams"
	"github.com/katallaxie/streams/msg"
)

var _ streams.Sink[any, any] = (*Sink[any, any])(nil)

// Sink is a noop sink.
type Sink[K, V any] struct {
	buf []msg.Message[K, V]
}

// NewSink is a noop sink constructor.
func NewSink[K, V any]() *Sink[K, V] {
	n := new(Sink[K, V])

	return n
}

// Write is a noop sink writer.
func (n *Sink[K, V]) Write(messages ...msg.Message[K, V]) error {
	n.buf = append(n.buf, messages...)

	return nil
}
