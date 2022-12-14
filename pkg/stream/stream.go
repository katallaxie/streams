package stream

type Message struct {
	Name string
}

// Sink ...
type Sink interface {
	Write(m *Message) error
}

// Stream ...
type Stream struct {
	in    chan *Message
	close chan bool
	err   chan error
}

// Close ...
func (s *Stream) Close() {
	close(s.in)
}

// Drain ...
func (s *Stream) Drain() {
	for _ = range s.in {
	}
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

	return &Stream{out, s.close, s.err}
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

	return &Stream{out, s.close, s.err}
}

// Branch ...
func (s *Stream) Branch(fns ...func(*Message) (bool, error)) []*Stream {
	streams := make([]*Stream, len(fns))

	for i := range fns {
		streams[i] = &Stream{make(chan *Message), s.close, s.err}
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
func (s *Stream) Sink(sink Sink) {
	for x := range s.in {
		err := sink.Write(x)
		if err != nil {
			s.Fail(err)
			return
		}
	}
}

// NewStream ...
func NewStream(in chan *Message) *Stream {
	stream := new(Stream)
	stream.in = in

	return stream
}
