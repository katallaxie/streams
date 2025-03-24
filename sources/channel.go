package sources

import (
	"github.com/katallaxie/streams"
)

// Source is the interface that represents a source of data.
type Source interface {
	streams.Streamable
}

var _ Source = (*ChanSource)(nil)

// ChanSource is a source that returns a channel of data.
type ChanSource struct {
	out chan any
}

// NewChanSource returns a new ChanSource.
func NewChanSource(out chan any) *ChanSource {
	return &ChanSource{out: out}
}

// Pipe pipes the output channel to the input channel.
func (s *ChanSource) Pipe(c streams.Connectable) streams.Connectable {
	streams.Pipe(s, c)
	return c
}

// Out returns the channel of data.
func (s *ChanSource) Out() <-chan any {
	return s.out
}
