package streams

import "context"

// Collector ...
type Collector interface {
	// Collect ...
	Collect(ch chan<- Metric)
}

// Metric
type Metric interface {
	// Write ...
	Write(m *Monitor) error
}

// Probe ...
type Probe[K, V any] interface {
	// Do ...
	Do(ctx context.Context, monitor Monitor) error

	Collector
}
