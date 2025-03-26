package streams

var (
	_ Streamable = (*PassThrough)(nil)
	_ Receivable = (*PassThrough)(nil)
)

// PassThrough passes through an incoming element.
type PassThrough struct {
	in  chan any
	out chan any
}

// NewPassThrough returns a new operator on pass-throughs.
func NewPassThrough() *PassThrough {
	t := &PassThrough{
		in:  make(chan any),
		out: make(chan any),
	}

	go t.attach()

	return t
}

// To streams data to the sink and waits for it to complete.
func (p *PassThrough) To(sink Sinkable) {
	p.stream(sink)
	sink.Wait()
}

// In returns the input channel.
func (p *PassThrough) In() chan<- any {
	return p.in
}

// Out returns the output channel.
func (p *PassThrough) Out() <-chan any {
	return p.out
}

// Pipe pipes the output channel to the input channel.
func (p *PassThrough) Pipe(c Operatable) Operatable {
	go p.stream(c)
	return c
}

func (m *PassThrough) stream(r Receivable) {
	for x := range m.out {
		r.In() <- x
	}

	close(r.In())
}

func (p *PassThrough) attach() {
	go func() {
		for x := range p.in {
			p.out <- x
		}

		close(p.out)
	}()
}
