package sinks

import (
	"github.com/katallaxie/streams"
)

var _ streams.Sinkable = (*ChanSink)(nil)

// ChanSink is a sink that receives data on a channel.
type ChanSink struct {
	in  chan<- any
	err chan error
}

// NewChanSink returns a new ChanSink.
func NewChanSink(in chan<- any) *ChanSink {
	return &ChanSink{in: in}
}

// Error returns the error.
func (s *ChanSink) Error() error {
	return <-s.err
}

// In returns the channel to send data to.
func (s *ChanSink) In() chan<- any {
	return s.in
}

// Wait waits for the sink to complete.
func (s *ChanSink) Wait() {
	// no-op
}
