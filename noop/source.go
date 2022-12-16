package noop

import (
	"github.com/ionos-cloud/streams/msg"
)

type noopSource struct {
	buf []msg.Message
	out chan msg.Message
}

// New ...
func NewSource(buf []msg.Message) *noopSource {
	n := new(noopSource)
	n.buf = buf
	n.out = make(chan msg.Message)

	return n
}

// Message ...
func (n *noopSource) Messages() chan msg.Message {
	out := make(chan msg.Message)

	go func(buf []msg.Message) {
		for _, msg := range n.buf {
			out <- msg
		}
	}(n.buf)

	return out
}

// Commit ...
func (n *noopSource) Commit(messages ...msg.Message) error {
	return nil
}

// Close ...
func (n *noopSource) Close() {
	close(n.out)
}
