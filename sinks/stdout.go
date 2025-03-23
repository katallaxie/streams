package sinks

import (
	"fmt"

	"github.com/katallaxie/streams"
)

// Stdout is a sink that writes data to stdout.
type Stdout struct {
	in   chan any
	done chan struct{}
}

var _ streams.Sinkable = (*ChanSink)(nil)

// NewStdout returns a new Stdout.
func NewStdout() *Stdout {
	out := &Stdout{
		in:   make(chan any),
		done: make(chan struct{}),
	}

	go out.process()

	return out
}

func (s *Stdout) process() {
	defer close(s.done)
	for elem := range s.in {
		fmt.Println(elem)
	}
}

// In returns the input channel of the StdoutSink connector.
func (s *Stdout) In() chan<- any {
	return s.in
}

// Wait waits for the sink to complete.
func (stdout *Stdout) Wait() {
	<-stdout.done
}
