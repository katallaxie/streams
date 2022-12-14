package stream

type Message struct {
	Name string
}

// Source ...
type Source interface {
	Messages() chan *Message
	Commit() error
}

// Sink ...
type Sink interface {
	Write(m *Message) error
}

// Stream ...
type Stream struct {
	in    chan *Message
	mark  chan *Message
	buf   chan *Message
	close chan bool
	err   chan error
}

// Close ...
func (s *Stream) Close() {
	close(s.in)
}

// Drain ...
func (s *Stream) Drain() {
	for range s.in {
	}
}

// Commit ...
func (s *Stream) Mark(m *Message) {
	s.mark <- m
}

// Fail ...
func (s *Stream) Fail(err error) error {
	s.Close()
	s.Drain()

	return err
}

// Filter ...
func (s *Stream) Filter(fn func(*Message) (bool, error)) *Stream {
	out := make(chan *Message)

	go func() {
		for x := range s.in {
			ok, err := fn(x)
			if err != nil {
				s.Fail(err)
				return
			}

			if ok {
				out <- x
			}
		}
		close(out)
	}()

	return &Stream{out, s.mark, s.buf, s.close, s.err}
}

// Map ...
func (s *Stream) Map(fn func(*Message) (*Message, error)) *Stream {
	out := make(chan *Message)

	go func() {
		for x := range s.in {
			x, err := fn(x)
			if err != nil {
				s.Fail(err)
				return
			}

			out <- x
		}
		close(out)
	}()

	return &Stream{out, s.mark, s.buf, s.close, s.err}
}

// Branch ...
func (s *Stream) Branch(fns ...func(*Message) (bool, error)) []*Stream {
	streams := make([]*Stream, len(fns))

	for i := range fns {
		streams[i] = &Stream{make(chan *Message), s.mark, s.buf, s.close, s.err}
	}

	go func() {
		for x := range s.in {
			for i, fn := range fns {
				ok, err := fn(x)
				if err != nil {
					s.Fail(err)
					return
				}

				if ok {
					streams[i].in <- x
				}
			}
		}

		for _, stream := range streams {
			close(stream.in)
		}
	}()

	return streams
}

// Sink ...
func (s *Stream) Sink(sink Sink) error {
	var err error

loop:
	for {
		select {
		case x, ok := <-s.in:
			if !ok {
				break loop
			}

			err := sink.Write(x)
			if err != nil {
				s.Fail(err)
			}

			s.mark <- x
		case err = <-s.err:
		}
	}

	return err
}

// NewStream ...
func NewStream(src Source, size int) *Stream {
	stream := new(Stream)
	stream.in = src.Messages()
	stream.mark = make(chan *Message)

	go func() {
		var count int
		var buf []*Message

		for m := range stream.mark {
			buf = append(buf, m)
			count++

			if len(stream.buf) <= size {
				continue
			}

			err := src.Commit()
			if err != nil {
				stream.Fail(err)
				return
			}

			buf = buf[:0]
			count = 0
		}
	}()

	return stream
}
