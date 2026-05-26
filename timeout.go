package streams

import (
	"time"
)

var (
	_ Streamable = (*TimeoutImpl)(nil)
	_ Receivable = (*TimeoutImpl)(nil)
)

// TimeoutImpl is an operator that closes the stream after a set amount of time.
type TimeoutImpl struct {
	dur time.Duration
	in  chan any
	out chan any
}

// Timeout returns a new timeout pipe.
func Timeout(dur time.Duration) *TimeoutImpl {
	return NewTimeout(dur)
}

// NewTimeout creates a new Timeout operator.
func NewTimeout(dur time.Duration) *TimeoutImpl {
	t := &TimeoutImpl{
		dur: dur,
		in:  make(chan any),
		out: make(chan any),
	}

	go t.attach()

	return t
}

// To streams data to the sink and waits for it to complete.
func (t *TimeoutImpl) To(sink Sinkable) error {
	t.stream(sink)

	err := sink.Wait()
	if err != nil {
		return err
	}

	return nil
}

// In returns the input channel.
func (t *TimeoutImpl) In() chan<- any {
	return t.in
}

// Out returns the output channel.
func (t *TimeoutImpl) Out() <-chan any {
	return t.out
}

// Pipe pipes the output channel to the input channel.
func (t *TimeoutImpl) Pipe(c Operatable) Operatable {
	go t.stream(c)
	return c
}

func (t *TimeoutImpl) stream(r Receivable) {
	elapsed := time.After(t.dur)
	defer close(t.out)

OUTTER:
	for {
		select {
		case v, ok := <-t.in:
			if !ok {
				break OUTTER
			}

			r.In() <- v

		case <-elapsed:
			break OUTTER
		}
	}

	close(r.In())
}

func (t *TimeoutImpl) attach() {
	go func() {
		for x := range t.in {
			_, ok := <-t.out
			if ok {
				t.out <- x
			}
		}

		close(t.out)
	}()
}
