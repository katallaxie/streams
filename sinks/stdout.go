package sinks

import (
	"fmt"

	"github.com/katallaxie/streams"
)

// Stdout is a sink that writes data to stdout.
type Stdout struct {
	in   chan any
	done chan struct{}
	err  chan error
}

var _ streams.Sinkable = (*ChanSink)(nil)

// NewStdout returns a new Stdout.
func NewStdout() *Stdout {
	out := &Stdout{
		in:   make(chan any),
		done: make(chan struct{}),
	}

	go out.attach()

	return out
}

// Error returns the error.
func (s *Stdout) Error() error {
	return <-s.err
}

func (s *Stdout) attach() {
	defer close(s.done)
	for elem := range s.in {
		fmt.Print(elem)
	}
}

// In returns the input channel of the StdoutSink connector.
func (s *Stdout) In() chan<- any {
	return s.in
}

// Wait waits for the sink to complete.
func (s *Stdout) Wait() {
	<-s.done
}
