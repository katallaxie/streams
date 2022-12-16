package streams

import (
	"log"
	"sync"

	"github.com/ionos-cloud/streams/msg"
)

// Source ...
type Source interface {
	Messages() chan *msg.Message
	Commit(...*msg.Message) error
}

// Sink ...
type Sink interface {
	Write(...*msg.Message) error
}

// Stream ...
type Stream struct {
	in    chan *msg.Message
	mark  chan *msg.Message
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

// Mark ...
func (s *Stream) Mark(m *msg.Message) {
	if s.mark == nil {
		return
	}

	s.mark <- m
}

// Fail ...
func (s *Stream) Fail(err error) {
	s.Close()
	s.Drain()

	s.err <- err
}

// Filter ...
func (s *Stream) Filter(fn func(*msg.Message) (bool, error)) *Stream {
	out := make(chan *msg.Message)

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

	return &Stream{out, s.mark, s.close, s.err}
}

// Map ...
func (s *Stream) Map(fn func(*msg.Message) (*msg.Message, error)) *Stream {
	out := make(chan *msg.Message)

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

	return &Stream{out, s.mark, s.close, s.err}
}

// Do ...
func (s *Stream) Do(fn func(*msg.Message)) *Stream {
	out := make(chan *msg.Message)

	go func() {
		for x := range s.in {
			fn(x)

			out <- x
		}
	}()

	return &Stream{out, s.mark, s.close, s.err}
}

// Branch ...
func (s *Stream) Branch(fns ...func(*msg.Message) (bool, error)) []*Stream {
	streams := make([]*Stream, len(fns))

	for i := range fns {
		streams[i] = &Stream{make(chan *msg.Message), s.mark, s.close, s.err}
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

// FanOut ...
func (s *Stream) FanOut(num int) []*Stream {
	streams := make([]*Stream, num)

	for i := range streams {
		streams[i] = &Stream{make(chan *msg.Message), s.mark, s.close, s.err}
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

// Print ...
func (s *Stream) Print() *Stream {
	out := make(chan *msg.Message)

	go func() {
		for x := range s.in {
			log.Printf(x.Name)

			out <- x
		}
		close(out)
	}()

	return &Stream{out, s.mark, s.close, s.err}
}

// Merge ...
func (s *Stream) Merge(streams ...*Stream) *Stream {
	out := make(chan *msg.Message)
	var closeOnce sync.Once

	for _, s := range streams {
		go func(c <-chan *msg.Message) {
			for x := range c {
				out <- x
			}

			closeOnce.Do(func() {
				close(out)
			})
		}(s.in)
	}

	return &Stream{out, s.mark, s.close, s.err}
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

			s.Mark(x)
		case err = <-s.err:
		}
	}

	return err
}

// NewStream ...
func NewStream(src Source, buffer int) *Stream {
	stream := new(Stream)
	stream.mark = make(chan *msg.Message)
	stream.in = src.Messages()

	go func() {
		var count int
		var buf []*msg.Message

		for m := range stream.mark {
			if m.Marked() {
				continue
			}

			buf = append(buf, m)
			count++

			m.Mark()

			if count <= buffer {
				continue
			}

			err := src.Commit(buf...)
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
