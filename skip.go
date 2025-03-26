package streams

var (
	_ Streamable = (*Reduce[any])(nil)
	_ Receivable = (*Reduce[any])(nil)
)

// Skip skips the first n elements.
type Skip struct {
	n   int
	in  chan any
	out chan any
}

// NewSkip returns a new operator on skips.
func NewSkip(n int) *Skip {
	t := &Skip{
		n:   n,
		in:  make(chan any),
		out: make(chan any),
	}

	go t.attach()

	return t
}

// To streams data to the sink and waits for it to complete.
func (s *Skip) To(sink Sinkable) {
	s.stream(sink)
	sink.Wait()
}

// In returns the input channel.
func (s *Skip) In() chan<- any {
	return s.in
}

// Out returns the output channel.
func (s *Skip) Out() <-chan any {
	return s.out
}

// Pipe pipes the output channel to the input channel.
func (s *Skip) Pipe(c Operatable) Operatable {
	go s.stream(c)
	return c
}

func (s *Skip) stream(recv Receivable) {
	for x := range s.out {
		recv.In() <- x
	}

	close(recv.In())
}

func (s *Skip) attach() {
	curr := s.n
	for x := range s.in {
		curr--
		if curr >= 0 {
			continue
		}
		s.out <- x
	}
	close(s.out)
}
