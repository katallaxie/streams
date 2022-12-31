package streams

import (
	"log"
	"sync"

	"github.com/ionos-cloud/streams/msg"
)

// Source is a source of messages.
type Source[K, V any] interface {
	Messages() chan msg.Message[K, V]
	Commit(...msg.Message[K, V]) error
}

// Sink is a sink of messages.
type Sink[K, V any] interface {
	Write(...msg.Message[K, V]) error
}

// Predicate is a function that returns true or false.
type Predicate[K, V any] func(msg.Message[K, V]) (bool, error)

// Stream is a stream of messages.
type Stream[K, V any] interface {
	Close()
	Do(fn func(msg.Message[K, V])) Stream[K, V]
	Drain()
	Fail(err error)
	FanOut(predicates ...Predicate[K, V]) []Stream[K, V]
	Filter(predicate Predicate[K, V]) Stream[K, V]
	Map(fn func(msg.Message[K, V]) (msg.Message[K, V], error)) Stream[K, V]
	Mark()
	Log() Stream[K, V]
	Sink(sink Sink[K, V]) error
}

// Close is a function that closes a stream.
func (s *StreamImpl[K, V]) Close() {
	close(s.in)
}

// Drain is a function that drains a stream.
func (s *StreamImpl[K, V]) Drain() {
	for range s.in {
	}
}

// Mark is a function that marks a message.
func (s *StreamImpl[K, V]) Mark(m msg.Message[K, V]) {
	if s.mark == nil {
		return
	}

	s.mark <- m
}

// Fail is a function that fails a stream
func (s *StreamImpl[K, V]) Fail(err error) {
	s.Close()
	s.Drain()

	s.err <- err
}

// Filter is a function that filters a stream.
func (s *StreamImpl[K, V]) Filter(fn Predicate[K, V]) *StreamImpl[K, V] {
	out := make(chan msg.Message[K, V])

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

	return &StreamImpl[K, V]{out, s.mark, s.close, s.err, s.opts}
}

// Map ...
func (s *StreamImpl[K, V]) Map(fn func(msg.Message[K, V]) (msg.Message[K, V], error)) *StreamImpl[K, V] {
	out := make(chan msg.Message[K, V])

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

	return &StreamImpl[K, V]{out, s.mark, s.close, s.err, s.opts}
}

// Do is a function that executes a function on a stream.
func (s *StreamImpl[K, V]) Do(fn func(msg.Message[K, V])) *StreamImpl[K, V] {
	out := make(chan msg.Message[K, V])

	go func() {
		for x := range s.in {
			fn(x)

			out <- x
		}
	}()

	return &StreamImpl[K, V]{out, s.mark, s.close, s.err, s.opts}
}

// Branch is branch a stream to multiple streams.
func (s *StreamImpl[K, V]) Branch(fns ...Predicate[K, V]) []*StreamImpl[K, V] {
	streams := make([]*StreamImpl[K, V], len(fns))

	for i := range fns {
		streams[i] = &StreamImpl[K, V]{make(chan msg.Message[K, V]), s.mark, s.close, s.err, s.opts}
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
func (s *StreamImpl[K, V]) FanOut(num int) []*StreamImpl[K, V] {
	streams := make([]*StreamImpl[K, V], num)

	for i := range streams {
		streams[i] = &StreamImpl[K, V]{make(chan msg.Message[K, V]), s.mark, s.close, s.err, s.opts}
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
func (s *StreamImpl[K, V]) Log() *StreamImpl[K, V] {
	out := make(chan msg.Message[K, V])

	go func() {
		for x := range s.in {
			log.Print(x.Key(), x.Value())

			out <- x
		}
		close(out)
	}()

	return &StreamImpl[K, V]{out, s.mark, s.close, s.err, s.opts}
}

// Merge ...
func (s *StreamImpl[K, V]) Merge(streams ...StreamImpl[K, V]) *StreamImpl[K, V] {
	out := make(chan msg.Message[K, V])
	var closeOnce sync.Once

	for _, s := range streams {
		go func(c <-chan msg.Message[K, V]) {
			for x := range c {
				out <- x
			}

			closeOnce.Do(func() {
				close(out)
			})
		}(s.in)
	}

	return &StreamImpl[K, V]{out, s.mark, s.close, s.err, s.opts}
}

// Sink is wire up a stream to a sink.
func (s *StreamImpl[K, V]) Sink(sink Sink[K, V]) error {
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
