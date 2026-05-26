package streams

import (
	"github.com/katallaxie/pkg/logx"
)

var (
	_ Streamable = (*LogImpl)(nil)
	_ Receivable = (*LogImpl)(nil)
)

// LogImpl passes through an incoming element.
type LogImpl struct {
	fn  logx.LogFunc
	in  chan any
	out chan any
}

// Log returns a new operator to log elements.
func Log(fn logx.LogFunc) *LogImpl {
	return NewLog(fn)
}

// NewLog returns a new operator to log elements.
func NewLog(fn logx.LogFunc) *LogImpl {
	l := &LogImpl{
		fn:  fn,
		in:  make(chan any),
		out: make(chan any),
	}

	go l.attach()

	return l
}

// To streams data to the sink and waits for it to complete.
func (l *LogImpl) To(sink Sinkable) error {
	l.stream(sink)

	err := sink.Wait()
	if err != nil {
		return err
	}

	return nil
}

// In returns the input channel.
func (l *LogImpl) In() chan<- any {
	return l.in
}

// Out returns the output channel.
func (l *LogImpl) Out() <-chan any {
	return l.out
}

// Pipe pipes the output channel to the input channel.
func (l *LogImpl) Pipe(c Operatable) Operatable {
	go l.stream(c)
	return c
}

func (l *LogImpl) stream(r Receivable) {
	go func() {
		for x := range l.in {
			l.fn.Printf("%v", x)
			r.In() <- x
		}

		close(l.out)
	}()
}

func (l *LogImpl) attach() {
	go func() {
		for x := range l.in {
			l.out <- x
		}

		close(l.out)
	}()
}
