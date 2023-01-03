package streams

import (
	"sync"
	"time"

	"github.com/ionos-cloud/streams/msg"
)

// Opts is a set of options for a stream.
type Opts struct {
	buffer   int
	nodeName string
	monitor  *Monitor
}

// Comfigure is a function that configures a stream.
func (o *Opts) Configure(opts ...Opt) {
	for _, opt := range opts {
		opt(o)
	}
}

// Opt is a function that configures a stream.
type Opt func(*Opts)

// WithBuffer configures the buffer size for a stream.
func WithBuffer(size int) Opt {
	return func(o *Opts) {
		o.buffer = size
	}
}

// WithNodeName configures the node name for a stream.
func WithNodeName(name string) Opt {
	return func(o *Opts) {
		o.nodeName = name
	}
}

// WithMonitor configures a statistics monitor.
func WithMonitor(m *Monitor) Opt {
	return func(o *Opts) {
		o.monitor = m
	}
}

// StreamImpl implements Stream.
type StreamImpl[K, V any] struct {
	in      chan msg.Message[K, V]
	mark    chan msg.Message[K, V]
	close   chan bool
	err     chan error
	metrics *metrics
	opts    *Opts

	Collector
}

type metrics struct {
	latency *latencyMetric
}

// NewStream from a source of messages.
func NewStream[K, V any](src Source[K, V], opts ...Opt) *StreamImpl[K, V] {
	options := new(Opts)
	options.Configure(opts...)

	out := make(chan msg.Message[K, V])

	stream := new(StreamImpl[K, V])
	stream.opts = options
	stream.mark = make(chan msg.Message[K, V])
	stream.in = out
	stream.err = make(chan error, 1)

	stream.metrics = new(metrics)
	stream.metrics.latency = newLatencyMetric(stream.opts.nodeName)

	go func() {
		for x := range src.Messages() {
			stream.metrics.latency.start()

			out <- x
		}

		close(out)
		close(stream.err)
	}()

	go func() {
		var count int
		var buf []msg.Message[K, V]

		for m := range stream.mark {
			if m.Marked() {
				continue
			}

			buf = append(buf, m)
			count++

			m.Mark()

			if count <= stream.opts.buffer {
				continue
			}

			err := src.Commit(buf...)
			if err != nil {
				stream.Fail(err)
				return
			}

			stream.metrics.latency.stop()

			if stream.opts.monitor != nil {
				stream.opts.monitor.Gather(stream)
			}

			buf = buf[:0]
			count = 0
		}
	}()

	return stream
}

type latencyMetric struct {
	value    float64
	nodeName string

	now time.Time

	Metric
	Collector

	sync.Mutex
}

// Collect is collecting metrics.
func (m *latencyMetric) Collect(ch chan<- Metric) {
	ch <- m
}

// Write is writing metrics to a channel.
func (m *latencyMetric) Write(monitor *Monitor) error {
	monitor.SetLatency(m.nodeName, m.value)

	return nil
}

func (m *latencyMetric) start() {
	m.Lock()
	defer m.Unlock()

	m.now = time.Now()
}

func (m *latencyMetric) stop() {
	m.Lock()
	defer m.Unlock()

	m.value = float64(time.Since(m.now).Microseconds())
}

func newLatencyMetric(nodeName string) *latencyMetric {
	return &latencyMetric{
		nodeName: nodeName,
		now:      time.Now(),
	}
}
