package streams

// MapFunc is a function that takes a key and a value and returns a new value.
type MapFunc[T, R any] func(T) R

var (
	_ Streamable = (*MapImpl[any, any])(nil)
	_ Receivable = (*MapImpl[any, any])(nil)
)

// MapImpl takes one element and produces a new element of the same type.
type MapImpl[T, R any] struct {
	fn  MapFunc[T, R]
	in  chan any
	out chan any
}

// Map returns a new operator on maps.
func Map[T, R any](fn MapFunc[T, R]) *MapImpl[T, R] {
	return NewMap(fn)
}

// NewMap returns a new operator on maps.
func NewMap[T, R any](fn MapFunc[T, R]) *MapImpl[T, R] {
	t := &MapImpl[T, R]{
		fn:  fn,
		in:  make(chan any),
		out: make(chan any),
	}

	go t.attach()

	return t
}

// To streams data to the sink and waits for it to complete.
func (m *MapImpl[T, R]) To(sink Sinkable) {
	m.stream(sink)
	sink.Wait()
}

// In returns the input channel.
func (m *MapImpl[T, R]) In() chan<- any {
	return m.in
}

// Out returns the output channel.
func (m *MapImpl[T, R]) Out() <-chan any {
	return m.out
}

// Pipe pipes the output channel to the input channel.
func (m *MapImpl[T, R]) Pipe(c Operatable) Operatable {
	go m.stream(c)
	return c
}

func (m *MapImpl[T, R]) stream(r Receivable) {
	for x := range m.out {
		r.In() <- x
	}

	close(r.In())
}

func (m *MapImpl[T, R]) attach() {
	for x := range m.in {
		x := m.fn(x.(T))
		m.out <- x
	}

	close(m.out)
}
