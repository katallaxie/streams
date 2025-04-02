package sinks

import (
	"github.com/katallaxie/pkg/fsmx"
	"github.com/katallaxie/streams"
)

// Store is a sink that stores all incoming data.
type Store struct {
	in      chan any
	store   fsmx.Store
	done    chan struct{}
	err     chan error
	mapFunc ActionFactory
}

var _ streams.Sinkable = (*ChanSink)(nil)

// ActionMapFunc is a function that maps any to an fsmx.Action.
type ActionFactory func(x any) fsmx.Action

// NewStore returns a new Store.
func NewStore(store fsmx.Store, fn ActionFactory) *Store {
	out := &Store{
		mapFunc: fn,
		store:   store,
		in:      make(chan any),
		done:    make(chan struct{}),
	}

	go out.attach()

	return out
}

// Error returns the error.
func (s *Store) Error() error {
	return <-s.err
}

func (s *Store) attach() {
	defer close(s.done)
	for x := range s.in {
		s.store.Dispatch(s.mapFunc(x))
	}
}

// In returns the input channel of the StdoutSink connector.
func (s *Store) In() chan<- any {
	return s.in
}

// Wait waits for the sink to complete.
func (stdout *Store) Wait() {
	<-stdout.done
}
