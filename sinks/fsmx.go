package sinks

import (
	"github.com/katallaxie/pkg/fsmx"
	"github.com/katallaxie/streams"
)

// FSMStoreSink is a sink that writes data to a FSM store.
type FSMStoreSink[S fsmx.State] struct {
	in    chan any
	store fsmx.Store[S]
	fn    ActionFactory
	done  chan struct{}
	err   chan error
}

// ActionFactory is a factory for creating actions.
type ActionFactory func(payload any) fsmx.Action

var _ streams.Sinkable = (*FSMStoreSink[fsmx.State])(nil)

// NewFSMStore return a new FSMStoreSink.
func NewFSMStore[S fsmx.State](store fsmx.Store[S], fn ActionFactory) *FSMStoreSink[S] {
	out := &FSMStoreSink[S]{
		store: store,
		fn:    fn,
		in:    make(chan any),
		done:  make(chan struct{}),
	}

	go out.attach()

	return out
}

// Error returns the error.
func (f *FSMStoreSink[S]) Error() error {
	return <-f.err
}

func (f *FSMStoreSink[S]) attach() {
	defer close(f.done)
	for elem := range f.in {
		action := f.fn(elem)
		f.store.Dispatch(action)
	}
}

// In returns the input channel of the StdoutSink connector.
func (f *FSMStoreSink[S]) In() chan<- any {
	return f.in
}

// Wait waits for the sink to complete.
func (f *FSMStoreSink[S]) Wait() {
	<-f.done
}
