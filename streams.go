package streams

import (
	"log"
	"sync"

	"github.com/ionos-cloud/streams/msg"
)

// Source is a source of messages.
type Source interface {
	Messages() chan msg.Message
	Commit(...msg.Message) error
}

// Sink is a sink of messages.
type Sink interface {
	Write(...msg.Message) error
}

// Predicate is a function that returns true or false.
type Predicate func(msg.Message) (bool, error)

// Stream is a stream of messages.
type Stream interface {
	Close()
	Do(fn func(msg.Message)) Stream
	Drain()
	Fail(err error)
	FanOut(predicates ...Predicate) []Stream
	Filter(fn func(msg.Message) (bool, error)) Stream
	Map(fn func(msg.Message) (msg.Message, error)) Stream
	Mark()
	Log() Stream
	Sink(sink Sink) error
}

// Close is a function that closes a stream.
func (s *StreamImpl) Close() {
	close(s.in)
}

// Drain is a function that drains a stream.
func (s *StreamImpl) Drain() {
	for range s.in {
	}
}

// Mark is a function that marks a message.
func (s *StreamImpl) Mark(m msg.Message) {
	if s.mark == nil {
		return
	}

	s.mark <- m
}

// Fail is a function that fails a stream
func (s *StreamImpl) Fail(err error) {
	s.Close()
	s.Drain()

	s.err <- err
}

// Filter is a function that filters a stream.
func (s *StreamImpl) Filter(fn func(msg.Message) (bool, error)) *StreamImpl {
	out := make(chan msg.Message)

	go func() {
		for x := range s.in {
			ok, err := fn(x)
			if err != nil {
				s.Fail(err)
				return
			}

			if ok {
				out <- x
			} else {
				s.Mark(x)
			}
		}
		close(out)
	}()

	return &StreamImpl{out, s.mark, s.close, s.err, s.opts}
}

// Map ...
func (s *StreamImpl) Map(fn func(msg.Message) (msg.Message, error)) *StreamImpl {
	out := make(chan msg.Message)

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

	return &StreamImpl{out, s.mark, s.close, s.err, s.opts}
}

// Do is a function that executes a function on a stream.
func (s *StreamImpl) Do(fn func(msg.Message)) *StreamImpl {
	out := make(chan msg.Message)

	go func() {
		for x := range s.in {
			fn(x)

			out <- x
		}
	}()

	return &StreamImpl{out, s.mark, s.close, s.err, s.opts}
}

// Branch is branch a stream to multiple streams.
func (s *StreamImpl) Branch(fns ...func(msg.Message) (bool, error)) []*StreamImpl {
	streams := make([]*StreamImpl, len(fns))

	for i := range fns {
		streams[i] = &StreamImpl{make(chan msg.Message), s.mark, s.close, s.err, s.opts}
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

// FanOut is fan out a stream to multiple streams.
func (s *StreamImpl) FanOut(num int) []*StreamImpl {
	streams := make([]*StreamImpl, num)

	for i := range streams {
		streams[i] = &StreamImpl{make(chan msg.Message), s.mark, s.close, s.err, s.opts}
	}

	go func() {
		for x := range s.in {
			for i := range streams {
				streams[i].in <- x
			}
		}

		for _, stream := range streams {
			close(stream.in)
		}
	}()

	return streams
}

// Log is logging the content of a stream.
func (s *StreamImpl) Log() *StreamImpl {
	out := make(chan msg.Message)

	go func() {
		for x := range s.in {
			log.Print(x)

			out <- x
		}
		close(out)
	}()

	return &StreamImpl{out, s.mark, s.close, s.err, s.opts}
}

// Merge ...
func (s *StreamImpl) Merge(streams ...StreamImpl) *StreamImpl {
	out := make(chan msg.Message)
	var closeOnce sync.Once

	for _, s := range streams {
		go func(c <-chan msg.Message) {
			for x := range c {
				out <- x
			}

			closeOnce.Do(func() {
				close(out)
			})
		}(s.in)
	}

	return &StreamImpl{out, s.mark, s.close, s.err, s.opts}
}

// Sink is wire up a stream to a sink.
func (s *StreamImpl) Sink(sink Sink) error {
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

			s.Mark(x)
		case err = <-s.err:
		}
	}

	return err
}
