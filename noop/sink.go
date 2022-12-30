package noop

import (
	"github.com/ionos-cloud/streams/msg"
)

type noopSink[K, V any] struct {
	buf []msg.Message[K, V]
}

// NewSink ...
func NewSink[K, V any]() *noopSink[K, V] {
	n := new(noopSink[K, V])

	return n
}

// Write ...
func (n *noopSink[K, V]) Write(messages ...msg.Message[K, V]) error {
	n.buf = append(n.buf, messages...)

	return nil
}
