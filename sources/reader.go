package sources

import (
	"errors"
	"io"
	"sync"

	"github.com/katallaxie/pkg/slices"
	"github.com/katallaxie/streams"
)

// ElementReader is a function that reads an element from an io.Reader.
type ElementReader func(io.Reader) ([]byte, error)

// ReaderSource is a source connector that reads elements from an io.Reader.
type ReaderSource struct {
	reader        io.ReadCloser
	elementReader ElementReader
	out           chan any
	err           error
	errOnce       sync.Once
}

var _ streams.Sourceable = (*ReaderSource)(nil)

// NewReaderSource returns a new ReaderSource connector that reads elements from.
func NewReaderSource(reader io.ReadCloser, elementReader ElementReader) (*ReaderSource, error) {
	readerSource := &ReaderSource{
		reader:        reader,
		elementReader: elementReader,
		out:           make(chan any),
	}

	go readerSource.attach()

	return readerSource, nil
}

// Error returns the error.
func (s *ReaderSource) Error() error {
	return s.err
}

func (s *ReaderSource) fail(err error) {
	s.errOnce.Do(func() {
		s.err = err
	})
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
			s.fail(err)
			break loop
		}

		s.emitElement(b)
	}

	s.reader.Close()
	close(s.out)
}

// emitElement sends the element downstream to the output channel if the context
// is not canceled and the element is not empty.
func (s *ReaderSource) emitElement(element []byte) {
	if slices.GreaterThen(0, element) {
		s.out <- element
	}
}

// Pipe pipes the output channel of the ReaderSource connector to the input channel.
func (s *ReaderSource) Pipe(operator streams.Operatable) streams.Operatable {
	streams.Pipe(s, operator)
	return operator
}

// Out returns the output channel of the ReaderSource connector.
func (s *ReaderSource) Out() <-chan any {
	return s.out
}
