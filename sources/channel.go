package sources

import (
	"github.com/katallaxie/streams"
)

var _ streams.Sourceable = (*ChanSourceImpl)(nil)

// ChanSource is a source that returns a channel of data.
type ChanSourceImpl struct {
	out chan any
}

// Channel returns a channel source.
func Channel(out chan any) *ChanSourceImpl {
	return NewChanSource(out)
}

// NewChanSource returns a new ChanSource.
func NewChanSource(out chan any) *ChanSourceImpl {
	return &ChanSourceImpl{out: out}
}

// Error returns the error.
func (s *ChanSourceImpl) Error() error {
	return nil // no-op
}

// Pipe pipes the output channel to the input channel.
func (s *ChanSourceImpl) Pipe(c streams.Operatable) streams.Operatable {
	streams.Pipe(s, c)
	return c
}

// Out returns the channel of data.
func (s *ChanSourceImpl) Out() <-chan any {
	return s.out
}
