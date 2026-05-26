package streams

var (
	_ Streamable = (*TakeTimpl)(nil)
	_ Receivable = (*TakeTimpl)(nil)
)

// TakeTimpl is a stream operator that takes a number of elements from the input channel.
type TakeTimpl struct {
	count int
	in    chan any
	out   chan any
}

// Take returns a new Take operator.
func Take(count int) *TakeTimpl {
	return NewTake(count)
}

// NewTake creates a new Take operator.
func NewTake(count int) *TakeTimpl {
	t := &TakeTimpl{
		count: count,
		in:    make(chan any),
		out:   make(chan any),
	}

	go t.attach()

	return t
}

// To streams data to the sink and waits for it to complete.
func (t *TakeTimpl) To(sink Sinkable) {
	t.stream(sink)
	sink.Wait()
}

// In returns the input channel.
func (t *TakeTimpl) In() chan<- any {
	return t.in
}

// Out returns the output channel.
func (t *TakeTimpl) Out() <-chan any {
	return t.out
}

// Pipe pipes the output channel to the input channel.
func (t *TakeTimpl) Pipe(c Operatable) Operatable {
	go t.stream(c)
	return c
}

func (t *TakeTimpl) stream(recv Receivable) {
	for x := range t.out {
		recv.In() <- x
	}

	close(recv.In())
}

func (t *TakeTimpl) attach() {
	go func() {
		for x := range t.in {
			if t.count > 0 {
				t.count--
				t.out <- x
			}

			if t.count == 0 {
				break
			}
		}

		close(t.out)
	}()
}
