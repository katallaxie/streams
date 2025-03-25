package streams

import (
	"sync"

	"github.com/katallaxie/pkg/slices"
)

// Pipe pipes the output channel to the input channel.
func Pipe(stream Streamable, rev Receivable) {
	go func() {
		for x := range stream.Out() {
			rev.In() <- x
		}

		close(rev.In())
	}()
}

// Streamable is a streamable interface.
type Streamable interface {
	// Out returns the output channel.
	Out() <-chan any
	// Pipe pipes the output channel to the input channel.
	Pipe(Connectable) Connectable
}

// Receivable is a receivable interface.
type Receivable interface {
	// In returns the input channel.
	In() chan<- any
}

// Sinkable is a sinkable interface.
type Sinkable interface {
	Receivable
	// Wait waits for the sink to complete.
	Wait()
	// Connect connects the sink to the source.
}

// Connectable is a connectable interface.
type Connectable interface {
	Streamable
	Receivable
	// To streams data to the sink and waits for it to complete.
	To(sink Sinkable)
}

// Split splits a stream in two based on a predicate.
func Split[T any](in Streamable, predicate FilterPredicate[T]) [2]Connectable {
	left := NewPassThrough()
	right := NewPassThrough()

	go func() {
		for x := range in.Out() {
			if predicate(x.(T)) {
				left.In() <- x
			} else {
				right.In() <- x
			}
		}
		close(left.In())
		close(right.In())
	}()

	return [...]Connectable{left, right}
}

// FanOut fans out a stream to multiple streams.
func FanOut(in Streamable, num int) []Connectable {
	out := make([]Connectable, num)

	slices.ForEach(func(o Connectable, i int) {
		out[i] = NewPassThrough()
	}, out...)

	go func() {
		for x := range in.Out() {
			for _, flow := range out {
				flow.In() <- x
			}
		}

		for _, flow := range out {
			close(flow.In())
		}
	}()

	return out
}

// Merge merges multiple streams into one.
func Merge(in ...Streamable) Connectable {
	merged := NewPassThrough()
	var wg sync.WaitGroup

	wg.Add(len(in))

	for _, out := range in {
		go func(in Streamable) {
			for element := range in.Out() {
				merged.In() <- element
			}

			wg.Done()
		}(out)
	}

	go func() {
		wg.Wait()
		close(merged.In())
	}()

	return merged
}

// Flatten creates a flatten stream.
func Flatten[T any]() Connectable {
	return NewFlatMap(func(element []T) []T { return element })
}
