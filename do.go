package streams

// DoFunc is a function that executes on the element.
type DoFunc[T any] func(T) error

var (
	_ Streamable = (*DoImpl[any])(nil)
	_ Receivable = (*DoImpl[any])(nil)
)

// DoImpl takes one element and executes a function on it.
type DoImpl[T any] struct {
	fn  DoFunc[T]
	in  chan any
	out chan any
}

// Do returns a new Do.
func Do[T any](fn DoFunc[T]) *DoImpl[T] {
	return NewDo(fn)
}

// NewDo creates a new Do.
func NewDo[T any](fn DoFunc[T]) *DoImpl[T] {
	t := &DoImpl[T]{
		fn:  fn,
		in:  make(chan any),
		out: make(chan any),
	}

	go t.attach()

	return t
}

// To streams data to the sink and waits for it to complete.
func (d *DoImpl[T]) To(sink Sinkable) error {
	d.stream(sink)

	err := sink.Wait()
	if err != nil {
		return err
	}

	return nil
}

// In returns the input channel.
func (d *DoImpl[T]) In() chan<- any {
	return d.in
}

// Out returns the output channel.
func (d *DoImpl[T]) Out() <-chan any {
	return d.out
}

// Pipe pipes the output channel to the input channel.
func (d *DoImpl[T]) Pipe(c Operatable) Operatable {
	go d.stream(c)
	return c
}

func (d *DoImpl[T]) stream(r Receivable) {
	for x := range d.out {
		r.In() <- x
	}

	close(r.In())
}

func (d *DoImpl[T]) attach() {
	for x := range d.in {
		err := d.fn(x.(T))
		if err != nil {
			break
		}
		d.out <- x
	}

	close(d.out)
}
