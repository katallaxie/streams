package sinks

import (
	"github.com/katallaxie/pkg/channels"
	"github.com/katallaxie/streams"
)

// Ignore represents a sink that ignores all incoming data.
type Ignore struct {
	in chan any
}

var _ streams.Sinkable = (*Ignore)(nil)

// NewIgnore returns a new sink that ignores the receivable data.
func NewIgnore() *Ignore {
	ignoreSink := &Ignore{
		in: make(chan any),
	}

	go ignoreSink.attach()

	return ignoreSink
}

func (i *Ignore) attach() {
	channels.Drain(i.in)
}

// In returns the input channel of a Ignore sink.
func (i *Ignore) In() chan<- any {
	return i.in
}

// Wait waits for the sink to complete.
func (i *Ignore) Wait() {
	// no-op
}
