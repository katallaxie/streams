package noop

import (
	"github.com/ionos-cloud/streams/msg"
)

// Sink ...
type Sink[K, V any] struct {
	buf []msg.Message[K, V]
}

// NewSink ...
func NewSink[K, V any]() *Sink[K, V] {
	n := new(Sink[K, V])

	return n
}

// Write ...
func (n *Sink[K, V]) Write(messages ...msg.Message[K, V]) error {
	n.buf = append(n.buf, messages...)

	return nil
}
