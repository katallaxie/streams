package streams

var (
	_ Streamable = (*SkipImpl)(nil)
	_ Receivable = (*SkipImpl)(nil)
)

// SkipImpl skips the first n elements.
type SkipImpl struct {
	n   int
	in  chan any
	out chan any
}

// Skip returns a new operator that skips the first n elements.
func Skip(n int) *SkipImpl {
	return NewSkip(n)
}

// NewSkip returns a new operator on skips.
func NewSkip(n int) *SkipImpl {
	t := &SkipImpl{
		n:   n,
		in:  make(chan any),
		out: make(chan any),
	}

	go t.attach()

	return t
}

// To streams data to the sink and waits for it to complete.
func (s *SkipImpl) To(sink Sinkable) error {
	s.stream(sink)

	err := sink.Wait()
	if err != nil {
		return err
	}

	return nil
}

// In returns the input channel.
func (s *SkipImpl) In() chan<- any {
	return s.in
}

// Out returns the output channel.
func (s *SkipImpl) Out() <-chan any {
	return s.out
}

// Pipe pipes the output channel to the input channel.
func (s *SkipImpl) Pipe(c Operatable) Operatable {
	go s.stream(c)
	return c
}

func (s *SkipImpl) stream(recv Receivable) {
	for x := range s.out {
		recv.In() <- x
	}

	close(recv.In())
}

func (s *SkipImpl) attach() {
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
