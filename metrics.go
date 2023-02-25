package streams

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// DefaultRegistry is a default prometheus registry.
	DefaultRegistry = NewRegistry()

	// DefaultRegisterer is a default prometheus registerer.
	DefaultRegisterer prometheus.Registerer = DefaultRegistry

	// DefaultGatherer is a default prometheus gatherer.
	DefaultGatherer prometheus.Gatherer = DefaultRegistry
)

// Registry is a prometheus registry.
type Registry struct {
	*prometheus.Registry
}

// NewRegistry is a constructor for Registry.
func NewRegistry() *Registry {
	r := prometheus.NewRegistry()

	r.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewBuildInfoCollector(),
	)

	return &Registry{Registry: r}
}

// Handler returns a HTTP handler for this registry. Should be registered at "/metrics".
func (r *Registry) Handler() http.Handler {
	return promhttp.InstrumentMetricHandler(r, promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
}

// DefaultMetrics is a set of default metrics.
var DefaultMetrics = NewMetrics()

// Metrics is a set of metrics.
type Metrics struct {
	latency *prometheus.GaugeVec
	count   *prometheus.CounterVec
}

// NewMetrics is a constructor for Metrics.
func NewMetrics() *Metrics {
	m := new(Metrics)

	m.latency = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "streams_latency",
			Help: "Time between sourcing and sinking a message.",
		},
		[]string{
			"streams_node",
		},
	)

	m.count = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "streams_count",
			Help: "Number of messages processed.",
		},
		[]string{
			"streams_node",
		},
	)

	return m
}

// Collect implements prometheus.Collector.
func (m *Metrics) Collect(ch chan<- prometheus.Metric) {
	m.latency.Collect(ch)
	m.count.Collect(ch)
}

// Describe implements prometheus.Collector.
func (m *Metrics) Describe(ch chan<- *prometheus.Desc) {
	m.latency.Describe(ch)
	m.count.Describe(ch)
}

// Monitor is a statistics monitor.
type Monitor struct {
	metrics *Metrics

	sync.Mutex
}

// Gatherer is a type that can gather metrics.
type Gatherer interface {
	// Gather ...
	Gather(collector Collector)
}

// NewMonitor is a constructor for Monitor.
func NewMonitor(metrics *Metrics) *Monitor {
	m := new(Monitor)
	m.metrics = metrics

	return m
}

// Gather is a method that gathers metrics.
func (m *Monitor) Gather(collector Collector) {
	m.Lock()
	defer m.Unlock()

	ch := make(chan Metric)
	defer func() { close(ch) }()

	go func() {
		for metric := range ch {
			_ = metric.Write(m)
		}
	}()

	collector.Collect(ch)
}

// SetLatency sets the latency metric.
func (m *Monitor) SetLatency(node string, latency float64) {
	m.metrics.latency.WithLabelValues(node).Set(latency)
}

// SetCount sets the count metric.
func (m *Monitor) SetCount(node string, count float64) {
	m.metrics.count.WithLabelValues(node).Add(count)
}
