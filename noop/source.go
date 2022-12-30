package noop

import (
	"github.com/ionos-cloud/streams/msg"
)

type noopSource[K, V any] struct {
	buf []msg.Message[K, V]
	out chan msg.Message[K, V]
}

// New ...
func NewSource[K, V any](buf []msg.Message[K, V]) *noopSource[K, V] {
	n := new(noopSource[K, V])
	n.buf = buf
	n.out = make(chan msg.Message[K, V])

	return n
}

// Message ...
func (n *noopSource[K, V]) Messages() chan msg.Message[K, V] {
	out := make(chan msg.Message[K, V])

	go func(buf []msg.Message[K, V]) {
		for _, msg := range n.buf {
			out <- msg
		}
	}(n.buf)

	return out
}

// Commit ...
func (n *noopSource[K, V]) Commit(messages ...msg.Message[K, V]) error {
	return nil
}

// Close ...
func (n *noopSource[K, V]) Close() {
	close(n.out)
}
