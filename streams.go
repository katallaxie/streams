package streams

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
