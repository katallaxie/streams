package sinks

import (
	"fmt"
	"io"

	"github.com/katallaxie/streams"
)

// Writer is a sink that writes data to an io.WriteCloser.
type Writer struct {
	writer io.WriteCloser
	in     chan any
	done   chan struct{}
	err    chan error
}

var _ streams.Sinkable = (*Writer)(nil)

// NewWriter creates a new WriterSink that writes data to the provided io.WriteCloser.
func NewWriterSink(writer io.WriteCloser) (*Writer, error) {
	w := &Writer{
		writer: writer,
		in:     make(chan any),
		done:   make(chan struct{}),
	}

	go w.attach()

	return w, nil
}

// Error returns the error.
func (w *Writer) Error() error {
	return <-w.err
}

// In returns the input channel of the WriterSink connector.
func (w *Writer) In() chan<- any {
	return w.in
}

// Wait waits for the sink to complete.
func (w *Writer) Wait() {
	<-w.done
}

func (w *Writer) attach() {
	defer close(w.done)

	for msg := range w.in {
		var bb []byte
		switch message := msg.(type) {
		case []byte:
			bb = message
		case string:
			bb = []byte(message)
		case fmt.Stringer:
			bb = []byte(message.String())
		default:
			continue
		}

		_, err := w.writer.Write(bb)
		if err != nil {
			break
		}
	}

	_ = w.writer.Close()
}
