package streams

// ReduceFunc combines the current element with the latest reduced value.
type ReduceFunc[T any] func(T, T) T

var (
	_ Streamable = (*Reduce[any])(nil)
	_ Receivable = (*Reduce[any])(nil)
)

// Reduce takes the current element and the latest reduced value and produces a new reduced value.
type Reduce[T any] struct {
	fn  ReduceFunc[T]
	in  chan any
	out chan any
}

// NewReduce returns a new operator on reduces.
func NewReduce[T any](fn ReduceFunc[T]) *Reduce[T] {
	t := &Reduce[T]{
		fn:  fn,
		in:  make(chan any),
		out: make(chan any),
	}

	go t.attach()

	return t
}

// To streams data to the sink and waits for it to complete.
func (r *Reduce[T]) To(sink Sinkable) {
	r.stream(sink)
	sink.Wait()
}

// In returns the input channel.
func (r *Reduce[T]) In() chan<- any {
	return r.in
}

// Out returns the output channel.
func (r *Reduce[T]) Out() <-chan any {
	return r.out
}

// Pipe pipes the output channel to the input channel.
func (r *Reduce[T]) Pipe(c Connectable) Connectable {
	go r.stream(c)
	return c
}

func (r *Reduce[T]) stream(recv Receivable) {
	for x := range r.out {
		recv.In() <- x
	}

	close(recv.In())
}

func (r *Reduce[T]) attach() {
	var curr any
	for x := range r.in {
		if curr == nil {
			curr = x
		} else {
			curr = r.fn(curr.(T), x.(T))
		}

		r.out <- curr
	}
	close(r.out)
}
