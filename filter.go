package streams

// FilterPredicate represents a filter predicate.
type FilterPredicate[T any] func(T) bool

var (
	_ Streamable = (*Filter[any])(nil)
	_ Receivable = (*Filter[any])(nil)
)

// Filter filters an incoming element using a filter predicate.
type Filter[T any] struct {
	fn  FilterPredicate[T]
	in  chan any
	out chan any
}

// NewFilter returns a new operator on filters.
func NewFilter[T any](fn FilterPredicate[T]) *Filter[T] {
	t := &Filter[T]{
		fn:  fn,
		in:  make(chan any),
		out: make(chan any),
	}

	go t.attach()

	return t
}

// To streams data to the sink and waits for it to complete.
func (f *Filter[T]) To(sink Sinkable) {
	f.stream(sink)
	sink.Wait()
}

// In returns the input channel.
func (f *Filter[T]) In() chan<- any {
	return f.in
}

// Out returns the output channel.
func (f *Filter[T]) Out() <-chan any {
	return f.out
}

// Pipe pipes the output channel to the input channel.
func (f *Filter[T]) Pipe(c Connectable) Connectable {
	go f.stream(c)
	return c
}

func (f *Filter[T]) stream(recv Receivable) {
	for x := range f.out {
		recv.In() <- x
	}

	close(recv.In())
}

func (f *Filter[T]) attach() {
	go func() {
		for x := range f.in {
			if f.fn(x.(T)) {
				f.out <- x
			}
		}

		close(f.out)
	}()
}
