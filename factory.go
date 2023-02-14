package streams

import (
	"sync"
	"time"

	"github.com/ionos-cloud/streams/msg"
)

// Opts is a set of options for a stream.
type Opts struct {
	buffer  int
	name    string
	monitor *Monitor
}

// Configure is a function that configures a stream.
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

// WithName configures the node name for a stream.
func WithName(name string) Opt {
	return func(o *Opts) {
		o.name = name
	}
}

// WithMonitor configures a statistics monitor.
func WithMonitor(m *Monitor) Opt {
	return func(o *Opts) {
		o.monitor = m
	}
}

// MessageChannel ...
type MessageChannel[K, V any] chan msg.Message[K, V]

// MessageReceiver ...
type MessageReceiver[K, V any] <-chan msg.Message[K, V]

// StreamImpl implements Stream.
type StreamImpl[K, V any] struct {
	in      MessageChannel[K, V]
	mark    MessageChannel[K, V]
	close   chan bool
	err     chan error
	metrics *metrics
	opts    *Opts

	topology Topology
	node     Node

	Collector
}

type metrics struct {
	latency *latencyMetric
	count   *countMetric
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

	node := NewNode("root")
	stream.node = node
	stream.topology = NewTopology(node)

	stream.metrics = new(metrics)
	stream.metrics.latency = newLatencyMetric(stream.opts.name)
	stream.metrics.count = newCountMetric(stream.opts.name)

	go func() {
		for x := range src.Messages() {
			stream.metrics.latency.start()

			out <- x
		}

		close(out)
		close(stream.err)
		close(stream.mark)
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
			stream.metrics.count.inc(len(buf))

			if stream.opts.monitor != nil {
				stream.opts.monitor.Gather(stream)
			}

			buf = buf[:0]
			count = 0
		}
	}()

	return stream
}

type countMetric struct {
	value    float64
	nodeName string

	Metric
	Collector

	sync.Mutex
}

// Collect is collecting metrics.
func (m *countMetric) Collect(ch chan<- Metric) {
	ch <- m
}

// Write is writing metrics to a channel.
func (m *countMetric) Write(monitor *Monitor) error {
	monitor.SetCount(m.nodeName, m.value)
	m.reset()

	return nil
}

func (m *countMetric) inc(count int) {
	m.Lock()
	defer m.Unlock()

	m.value += float64(int64(count))
}

func (m *countMetric) reset() {
	m.Lock()
	defer m.Unlock()

	m.value = 0
}

func newCountMetric(nodeName string) *countMetric {
	return &countMetric{
		nodeName: nodeName,
	}
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
