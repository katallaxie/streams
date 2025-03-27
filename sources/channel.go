package sources

import (
	"github.com/katallaxie/streams"
)

var _ streams.Sourceable = (*ChanSource)(nil)

// ChanSource is a source that returns a channel of data.
type ChanSource struct {
	out chan any
}

// NewChanSource returns a new ChanSource.
func NewChanSource(out chan any) *ChanSource {
	return &ChanSource{out: out}
}

// Pipe pipes the output channel to the input channel.
func (s *ChanSource) Pipe(c streams.Operatable) streams.Operatable {
	streams.Pipe(s, c)
	return c
}

// Out returns the channel of data.
func (s *ChanSource) Out() <-chan any {
	return s.out
}
