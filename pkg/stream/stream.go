package stream

type Message struct{}

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

// Fail ...
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
	}()

	return &Stream{out, s.close, s.err}
}

// NewStream ...
func NewStream(in chan *Message) *Stream {
	stream := new(Stream)
	stream.in = in

	return stream
}

// // Map ...
// func (in Stream) Map(fn func(*Message) *Message) Stream {
// 	out := make(Stream)

// 	go func() {
// 		for x := range in {
// 			out <- fn(x)
// 		}
// 	}()

// 	return out
// }

// // Do ...
// func (in Stream) Do(fn func(*Message)) Stream {
// 	out := make(Stream)

// 	for x := range in {
// 		fn(x)

// 		out <- x
// 	}

// 	return out
// }
