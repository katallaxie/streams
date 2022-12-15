package noop

import (
	"github.com/ionos-cloud/streams/msg"
)

type noopSink struct {
	buf []*msg.Message
}

// NewSink ...
func NewSink() *noopSink {
	n := new(noopSink)

	return n
}

// Write ...
func (n *noopSink) Write(messages ...*msg.Message) error {
	n.buf = append(n.buf, messages...)

	return nil
}
