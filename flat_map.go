package streams

// FlatMap takes one element and produces a new element of the same type.
type FlatMapFunc[T, R any] func(T) []R

var (
	_ Streamable = (*FlatMap[any, any])(nil)
	_ Receivable = (*FlatMap[any, any])(nil)
)

// FlatMap takes one element and produces a new element of the same type.
type FlatMap[T, R any] struct {
	fn  FlatMapFunc[T, R]
	in  chan any
	out chan any
}

// NewFlatMap returns a new operator on maps.
func NewFlatMap[T, R any](fn FlatMapFunc[T, R]) *FlatMap[T, R] {
	t := &FlatMap[T, R]{
		fn:  fn,
		in:  make(chan any),
		out: make(chan any),
	}

	go t.attach()

	return t
}

// To streams data to the sink and waits for it to complete.
func (f *FlatMap[T, R]) To(sink Sinkable) {
	f.stream(sink)
	sink.Wait()
}

// In returns the input channel.
func (f *FlatMap[T, R]) In() chan<- any {
	return f.in
}

// Out returns the output channel.
func (f *FlatMap[T, R]) Out() <-chan any {
	return f.out
}

// Pipe pipes the output channel to the input channel.
func (f *FlatMap[T, R]) Pipe(c Operatable) Operatable {
	go f.stream(c)
	return c
}

func (f *FlatMap[T, R]) stream(r Receivable) {
	for x := range f.out {
		r.In() <- x
	}

	close(r.In())
}

func (f *FlatMap[T, R]) attach() {
	for x := range f.in {
		x := f.fn(x.(T))
		for _, y := range x {
			f.out <- y
		}
	}

	close(f.out)
}
