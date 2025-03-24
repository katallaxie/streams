package streams

import (
	"github.com/katallaxie/pkg/logx"
)

var (
	_ Streamable = (*Log)(nil)
	_ Receivable = (*Log)(nil)
)

// Log passes through an incoming element.
type Log struct {
	fn  logx.LogFunc
	in  chan any
	out chan any
}

// NewLog returns a new operator to log elements.
func NewLog(fn logx.LogFunc) *Log {
	l := &Log{
		fn:  fn,
		in:  make(chan any),
		out: make(chan any),
	}

	go l.attach()

	return l
}

// To streams data to the sink and waits for it to complete.
func (l *Log) To(sink Sinkable) {
	l.stream(sink)
	sink.Wait()
}

// In returns the input channel.
func (l *Log) In() chan<- any {
	return l.in
}

// Out returns the output channel.
func (l *Log) Out() <-chan any {
	return l.out
}

// Pipe pipes the output channel to the input channel.
func (l *Log) Pipe(c Connectable) Connectable {
	go l.stream(c)
	return c
}

func (l *Log) stream(r Receivable) {
	go func() {
		for x := range l.in {
			l.fn.Printf("%v", x)
			r.In() <- x
		}

		close(l.out)
	}()
}

func (l *Log) attach() {
	go func() {
		for x := range l.in {
			l.out <- x
		}

		close(l.out)
	}()
}
