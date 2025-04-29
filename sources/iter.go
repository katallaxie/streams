package sources

import (
	"iter"

	"github.com/katallaxie/streams"
)

// SeqSource is a source that iterates over an iterable.
type SeqSource[I any] struct {
	seq iter.Seq[I]
	out chan any
}

var _ streams.Sourceable = (*SeqSource[any])(nil)

// NewSeqSource returns a new SeqSource.
func NewSeqSource[I any](seq iter.Seq[I]) (*SeqSource[I], error) {
	seqSource := &SeqSource[I]{
		seq: seq,
		out: make(chan any),
	}

	go seqSource.attach()

	return seqSource, nil
}

// Error returns the error.
func (s *SeqSource[I]) Error() error {
	return nil
}

func (s *SeqSource[I]) attach() {
	for e := range s.seq {
		s.out <- e
	}
	close(s.out)
}

// Pipe pipes the output channel of the ReaderSource connector to the input channel.
func (s *SeqSource[I]) Pipe(operator streams.Operatable) streams.Operatable {
	streams.Pipe(s, operator)
	return operator
}

// Out returns the output channel of the ReaderSource connector.
func (s *SeqSource[I]) Out() <-chan any {
	return s.out
}
