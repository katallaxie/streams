package streams

// MapFunction is a function that takes a key and a value and returns a new value.

type MapFunction[T, R any] func(T) R

var (
	_ Streamable = (*Map[any, any])(nil)
	_ Receivable = (*Map[any, any])(nil)
)

// Map takes one element and produces a new element of the same type.
type Map[T, R any] struct {
	fn  MapFunction[T, R]
	in  chan any
	out chan any
}

// NewMap returns a new operator on maps.
func NewMap[T, R any](fn MapFunction[T, R]) *Map[T, R] {
	t := &Map[T, R]{
		fn:  fn,
		in:  make(chan any),
		out: make(chan any),
	}

	t.attach()

	return t
}

// To streams data to the sink and waits for it to complete.
func (m *Map[T, R]) To(sink Sinkable) {
	m.stream(sink)
	sink.Wait()
}

// In returns the input channel.
func (m *Map[T, R]) In() chan<- any {
	return m.in
}

// Out returns the output channel.
func (m *Map[T, R]) Out() <-chan any {
	return m.out
}

// Pipe pipes the output channel to the input channel.
func (m *Map[T, R]) Pipe(c Connectable) Connectable {
	go m.stream(c)
	return c
}

func (m *Map[T, R]) stream(r Receivable) {
	for x := range m.out {
		r.In() <- x
	}

	close(r.In())
}

func (m *Map[T, R]) attach() {
	go func() {
		for x := range m.in {
			x := m.fn(x.(T))
			m.out <- x
		}

		close(m.out)
	}()
}
