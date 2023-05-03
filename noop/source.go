package noop

import (
	"github.com/katallaxie/streams"
	"github.com/katallaxie/streams/msg"
)

var _ streams.Source[any, any] = (*Source[any, any])(nil)

// Source ...
type Source[K, V any] struct {
	buf []msg.Message[K, V]
	out chan msg.Message[K, V]
}

// NewSource is a noop source constructor.
func NewSource[K, V any](buf []msg.Message[K, V]) *Source[K, V] {
	n := new(Source[K, V])
	n.buf = buf
	n.out = make(chan msg.Message[K, V])

	return n
}

// Messages is returning a channel of messages.
func (n *Source[K, V]) Messages() chan msg.Message[K, V] {
	out := make(chan msg.Message[K, V])

	go func(buf []msg.Message[K, V]) {
		for _, msg := range n.buf {
			out <- msg
		}
	}(n.buf)

	return out
}

// Commit is a noop source commit.
func (n *Source[K, V]) Commit(messages ...msg.Message[K, V]) error {
	return nil
}

// Close is a noop source closer.
func (n *Source[K, V]) Close() {
	close(n.out)
}

// Error is a noop source error.
func (n *Source[K, V]) Error() error {
	return nil
}
