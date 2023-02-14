package noop

import (
	"github.com/ionos-cloud/streams/msg"
)

// Source ...
type Source[K, V any] struct {
	buf []msg.Message[K, V]
	out chan msg.Message[K, V]
}

// NewSource ...
func NewSource[K, V any](buf []msg.Message[K, V]) *Source[K, V] {
	n := new(Source[K, V])
	n.buf = buf
	n.out = make(chan msg.Message[K, V])

	return n
}

// Message ...
func (n *Source[K, V]) Messages() chan msg.Message[K, V] {
	out := make(chan msg.Message[K, V])

	go func(buf []msg.Message[K, V]) {
		for _, msg := range n.buf {
			out <- msg
		}
	}(n.buf)

	return out
}

// Commit ...
func (n *Source[K, V]) Commit(messages ...msg.Message[K, V]) error {
	return nil
}

// Close ...
func (n *Source[K, V]) Close() {
	close(n.out)
}
