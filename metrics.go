package streams

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	DefaultRegistry                         = NewRegistry()
	DefaultRegisterer prometheus.Registerer = DefaultRegistry
	DefaultGatherer   prometheus.Gatherer   = DefaultRegistry
)

// Registry ...
type Registry struct {
	*prometheus.Registry
}

// NewRegistry ...
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

var (
	// DefaultMetrics ...
	DefaultMetrics = NewMetrics()
)

// Metrics ...
type Metrics struct {
}

// NewMetrics ...
func NewMetrics() *Metrics {
	m := new(Metrics)

	return m
}

// Collect ...
func (m *Metrics) Collect(ch chan<- prometheus.Metric) {

}

// Describe ...
func (m *Metrics) Describe(ch chan<- *prometheus.Desc) {

}

// Monitor ...
type Monitor struct {
	metrics *Metrics

	sync.Mutex
}

// Gatherer ...
type Gatherer interface {
	// Gather ...
	Gather(collector Collector)
}

// NewMonitor ...
func NewMonitor(metrics *Metrics) *Monitor {
	m := new(Monitor)
	m.metrics = metrics

	return m
}

// Gather ...
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
