package sources

import (
	"errors"
	"io"

	"github.com/katallaxie/streams"
)

// ElementReader is a function that reads an element from an io.Reader.
type ElementReader func(io.Reader) ([]byte, error)

// ReaderSource is a source connector that reads elements from an io.Reader.
type ReaderSource struct {
	reader        io.ReadCloser
	elementReader ElementReader
	out           chan any
}

var _ Source = (*ReaderSource)(nil)

// NewReaderSource returns a new ReaderSource connector that reads elements from
func NewReaderSource(reader io.ReadCloser, elementReader ElementReader) (*ReaderSource, error) {
	readerSource := &ReaderSource{
		reader:        reader,
		elementReader: elementReader,
		out:           make(chan any),
	}

	go readerSource.attach()

	return readerSource, nil
}

func (s *ReaderSource) attach() {
loop:
	for {
		b, err := s.elementReader(s.reader)
		if errors.Is(err, io.EOF) {
			s.emitElement(b)
			break loop
		}

		if err != nil {
			break loop
		}

		s.emitElement(b)
	}

	s.reader.Close()
	close(s.out)
}

// emitElement sends the element downstream to the output channel if the context
// is not canceled and the element is not empty.
func (r *ReaderSource) emitElement(element []byte) {
	if len(element) > 0 {
		r.out <- element
	}
}

// Pipe pipes the output channel of the ReaderSource connector to the input channel
func (s *ReaderSource) Pipe(operator streams.Connectable) streams.Connectable {
	streams.Pipe(s, operator)
	return operator
}

// Out returns the output channel of the ReaderSource connector.
func (s *ReaderSource) Out() <-chan any {
	return s.out
}
