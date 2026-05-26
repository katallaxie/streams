package streams

var (
	_ Streamable = (*Take)(nil)
	_ Receivable = (*Take)(nil)
)

// Take is a stream operator that takes a number of elements from the input channel.
type Take struct {
	count int
	in    chan any
	out   chan any
}

// NewTake creates a new Take operator.
func NewTake(count int) *Take {
	t := &Take{
		count: count,
		in:    make(chan any),
		out:   make(chan any),
	}

	go t.attach()

	return t
}

// To streams data to the sink and waits for it to complete.
func (f *Take) To(sink Sinkable) {
	f.stream(sink)
	sink.Wait()
}

// In returns the input channel.
func (f *Take) In() chan<- any {
	return f.in
}

// Out returns the output channel.
func (f *Take) Out() <-chan any {
	return f.out
}

// Pipe pipes the output channel to the input channel.
func (f *Take) Pipe(c Operatable) Operatable {
	go f.stream(c)
	return c
}

func (f *Take) stream(recv Receivable) {
	for x := range f.out {
		recv.In() <- x
	}

	close(recv.In())
}

func (f *Take) attach() {
	go func() {
		for x := range f.in {
			if f.count > 0 {
				f.count--
				f.out <- x
			}

			if f.count == 0 {
				break
			}
		}

		close(f.out)
	}()
}
