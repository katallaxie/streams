package streams

// DoFunc is a function that executes on the element.
type DoFunc[T any] func(T) error

var (
	_ Streamable = (*Map[any, any])(nil)
	_ Receivable = (*Map[any, any])(nil)
)

// Do takes one element and executes a function on it.
type Do[T any] struct {
	fn  DoFunc[T]
	in  chan any
	out chan any
}

// NewDo creates a new Do.
func NewDo[T any](fn DoFunc[T]) *Do[T] {
	t := &Do[T]{
		fn:  fn,
		in:  make(chan any),
		out: make(chan any),
	}

	go t.attach()

	return t
}

// To streams data to the sink and waits for it to complete.
func (d *Do[T]) To(sink Sinkable) {
	d.stream(sink)
	sink.Wait()
}

// In returns the input channel.
func (d *Do[T]) In() chan<- any {
	return d.in
}

// Out returns the output channel.
func (d *Do[T]) Out() <-chan any {
	return d.out
}

// Pipe pipes the output channel to the input channel.
func (d *Do[T]) Pipe(c Operatable) Operatable {
	go d.stream(c)
	return c
}

func (d *Do[T]) stream(r Receivable) {
	for x := range d.out {
		r.In() <- x
	}

	close(r.In())
}

func (d *Do[T]) attach() {
	for x := range d.in {
		err := d.fn(x.(T))
		if err != nil {
			break
		}
		d.out <- x
	}

	close(d.out)
}
