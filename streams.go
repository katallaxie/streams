package streams

import (
	"math"
	"sync"
	"time"

	"github.com/katallaxie/streams/msg"
)

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

// Source is a source of messages.
type Source[K, V any] interface {
	// Messages returns a channel of messages.
	Messages() chan msg.Message[K, V]
	// Commit commits a message.
	Commit(...msg.Message[K, V]) error
	// Error returns an error.
	Error() error
}

// Sink is a sink of messages.
type Sink[K, V any] interface {
	Write(...msg.Message[K, V]) error
}

// Predicate is a function that returns true or false.
type Predicate[K, V any] func(msg.Message[K, V]) (bool, error)

// Key is a message key.
type Key interface {
	int | ~string | []byte
}

// Value is a message value.
type Value interface {
	int | ~string | []byte
}

// Messages is a channel of messages.
type Messages[K, V any] chan msg.Message[K, V]

// Stream is a stream of messages.
type Stream[K Key, V Value] interface {
	// Close closes a stream.
	Close()

	// Do executes a function on a stream.
	Do(name string, fn func(msg.Message[K, V])) Stream[K, V]

	// Drain drains a stream.
	Drain()

	// FanOut splits a stream into multiple streams.
	FanOut(name string, predicates ...Predicate[K, V]) []Stream[K, V]

	// Filter filters a stream.
	Filter(name string, predicate Predicate[K, V]) Stream[K, V]

	// Map maps a stream.
	Map(name string, fn func(msg.Message[K, V]) (msg.Message[K, V], error)) Stream[K, V]

	// Mark marks a message.
	Mark() Stream[K, V]

	// Log logs a message.
	Log(name string) Stream[K, V]

	// Sink sends messages to a sink.
	Sink(name string, sink Sink[K, V])

	// Errors returns the first error.
	Error() error
}

// Unimplemented ...
type Unimplemented[K Key, V Value] struct{}

// Close is a function that closes a stream.
func (u *Unimplemented[K, V]) Close() {}

// Drain is a function that drains a stream.
func (u *Unimplemented[K, V]) Drain() {}

// Mark is a function that marks a message.
func (u *Unimplemented[K, V]) Mark() Stream[K, V] {
	return u
}

// Log is a function that logs a message.
func (u *Unimplemented[K, V]) Log(name string) Stream[K, V] {
	return u
}

// Sink is a function that sends messages to a sink.
func (u *Unimplemented[K, V]) Sink(name string, sink Sink[K, V]) {}

// Error is a function that returns the first error.
func (u *Unimplemented[K, V]) Error() error {
	return nil
}

// Filter is a function that filters a stream.
func (u *Unimplemented[K, V]) Filter(name string, predicate Predicate[K, V]) Stream[K, V] {
	return u
}

// Do is a function that executes a function on a stream.
func (u *Unimplemented[K, V]) Do(name string, fn func(msg.Message[K, V])) Stream[K, V] {
	return u
}

// FanOut is a function that splits a stream into multiple streams.
func (u *Unimplemented[K, V]) FanOut(name string, predicates ...Predicate[K, V]) []Stream[K, V] {
	return nil
}

// Map is a function that maps a stream.
func (u *Unimplemented[K, V]) Map(name string, fn func(msg.Message[K, V]) (msg.Message[K, V], error)) Stream[K, V] {
	return u
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
func (s *StreamImpl[K, V]) Mark(name string) *StreamImpl[K, V] {
	out := make(chan msg.Message[K, V])

	node := NewNode(name)
	s.node.AddChild(node)

	go func() {
		for x := range s.in {
			x.Mark()
			out <- x
		}
		close(out)
	}()

	return &StreamImpl[K, V]{out, s.src, s.close, s.err, s.metrics, s.opts, s.topology, node, s.Collector}
}

// Fail is a function that fails a stream
func (s *StreamImpl[K, V]) Fail(err error) {
	s.Drain()
	s.Close()

	s.error().Printf("%v", err)

	s.err <- err
}

// Error is a function that returns the error of a stream.
func (s *StreamImpl[K, V]) Error() error {
	return <-s.err
}

// Filter is a function that filters a stream.
func (s *StreamImpl[K, V]) Filter(name string, fn Predicate[K, V]) *StreamImpl[K, V] {
	out := make(chan msg.Message[K, V])

	node := NewNode(name)
	s.node.AddChild(node)

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
				x.Mark()
			}
		}
		close(out)
	}()

	return &StreamImpl[K, V]{out, s.src, s.close, s.err, s.metrics, s.opts, s.topology, node, s.Collector}
}

// Map is a function that maps a stream.
func (s *StreamImpl[K, V]) Map(name string, fn func(msg.Message[K, V]) (msg.Message[K, V], error)) *StreamImpl[K, V] {
	out := make(chan msg.Message[K, V])

	node := NewNode(name)
	s.node.AddChild(node)

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

	return &StreamImpl[K, V]{out, s.src, s.close, s.err, s.metrics, s.opts, s.topology, node, s.Collector}
}

// Do is a function that executes a function on a stream.
func (s *StreamImpl[K, V]) Do(name string, fn func(msg.Message[K, V])) *StreamImpl[K, V] {
	out := make(chan msg.Message[K, V])

	node := NewNode(name)
	s.node.AddChild(node)

	go func() {
		for x := range s.in {
			fn(x)

			out <- x
		}
	}()

	return &StreamImpl[K, V]{out, s.src, s.close, s.err, s.metrics, s.opts, s.topology, node, s.Collector}
}

// Branch is branch a stream to multiple streams.
func (s *StreamImpl[K, V]) Branch(name string, fns ...Predicate[K, V]) []*StreamImpl[K, V] {
	streams := make([]*StreamImpl[K, V], len(fns))

	for i := range fns {
		node := NewNode(name)
		s.node.AddChild(node)

		streams[i] = &StreamImpl[K, V]{make(chan msg.Message[K, V]), s.src, s.close, s.err, s.metrics, s.opts, s.topology, node, s.Collector}
	}

	go func() {
		for x := range s.in {
			mark := true

			for i, fn := range fns {
				ok, err := fn(x)
				if err != nil {
					s.Fail(err)
					return
				}

				if ok {
					streams[i].in <- x
					mark = false
					break
				}
			}

			if mark {
				x.Mark()
			}
		}

		for _, stream := range streams {
			close(stream.in)
		}
	}()

	return streams
}

// FanOut is fan out a stream to multiple streams.
func (s *StreamImpl[K, V]) FanOut(name string, num int) []*StreamImpl[K, V] {
	streams := make([]*StreamImpl[K, V], num)

	for i := range streams {
		node := NewNode(name)
		s.node.AddChild(node)

		streams[i] = &StreamImpl[K, V]{make(chan msg.Message[K, V]), s.src, s.close, s.err, s.metrics, s.opts, s.topology, node, s.Collector}
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
func (s *StreamImpl[K, V]) Log(name string) *StreamImpl[K, V] {
	out := make(chan msg.Message[K, V])

	node := NewNode(name)
	s.node.AddChild(node)

	go func() {
		for x := range s.in {
			s.log().Printf(name, "key", x.Key(), "partition", x.Partition(), "offset", x.Offset())

			out <- x
		}

		close(out)
	}()

	return &StreamImpl[K, V]{out, s.src, s.close, s.err, s.metrics, s.opts, s.topology, node, s.Collector}
}

// Merge is merge multiple streams into one.
func (s *StreamImpl[K, V]) Merge(name string, streams ...StreamImpl[K, V]) *StreamImpl[K, V] {
	out := make(chan msg.Message[K, V])
	var closeOnce sync.Once

	node := NewNode(name)
	s.node.AddChild(node)

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

	return &StreamImpl[K, V]{out, s.src, s.close, s.err, s.metrics, s.opts, s.topology, node, s.Collector}
}

// Sink is wire up a stream to a sink.
// nolint: gocyclo
func (s *StreamImpl[K, V]) Sink(name string, sink Sink[K, V]) {
	node := NewNode(name)
	s.node.AddChild(node)

	go func(c <-chan msg.Message[K, V]) {
		var buf []msg.Message[K, V]
		var count int

		ticker := time.NewTicker(math.MaxInt32 * time.Second)
		defer ticker.Stop()

		if s.opts.timeout > 0 {
			ticker.Reset(s.opts.timeout)
		} else {
			ticker.Stop()
		}

	LOOP:
		for {
			select {
			case <-ticker.C:
				err := sink.Write(buf...)
				if err != nil {
					s.Fail(err)
					return
				}

				for _, m := range buf {
					m.Mark()
				}

				buf = buf[:0]
				count = 0

				ticker.Reset(s.opts.timeout)
			case m, ok := <-c:
				if !ok {
					break LOOP
				}

				buf = append(buf, m)
				count++

				if s.opts.timeout > 0 {
					continue
				}

				if count <= s.opts.buffer {
					continue
				}

				err := sink.Write(buf...)
				if err != nil {
					s.Fail(err)
					break LOOP
				}

				for _, m := range buf {
					m.Mark()
				}

				buf = buf[:0]
				count = 0
			}
		}

		close(s.err)
	}(s.in)
}

// Collect is collect the content of a stream.
func (s *StreamImpl[K, V]) Collect(ch chan<- Metric) {
	s.metrics.count.Collect(ch)
}
