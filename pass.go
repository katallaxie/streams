package streams

var (
	_ Streamable = (*PassThroughImpl)(nil)
	_ Receivable = (*PassThroughImpl)(nil)
)

// PassThroughImpl passes through an incoming element.
type PassThroughImpl struct {
	in  chan any
	out chan any
}

// PassThrough returns a new operator on pass-throughs.
func PassThrough() *PassThroughImpl {
	return NewPassThrough()
}

// NewPassThrough returns a new operator on pass-throughs.
func NewPassThrough() *PassThroughImpl {
	t := &PassThroughImpl{
		in:  make(chan any),
		out: make(chan any),
	}

	go t.attach()

	return t
}

// To streams data to the sink and waits for it to complete.
func (p *PassThroughImpl) To(sink Sinkable) {
	p.stream(sink)
	sink.Wait()
}

// In returns the input channel.
func (p *PassThroughImpl) In() chan<- any {
	return p.in
}

// Out returns the output channel.
func (p *PassThroughImpl) Out() <-chan any {
	return p.out
}

// Pipe pipes the output channel to the input channel.
func (p *PassThroughImpl) Pipe(c Operatable) Operatable {
	go p.stream(c)
	return c
}

func (p *PassThroughImpl) stream(r Receivable) {
	for x := range p.out {
		r.In() <- x
	}

	close(r.In())
}

func (p *PassThroughImpl) attach() {
	go func() {
		for x := range p.in {
			p.out <- x
		}

		close(p.out)
	}()
}
