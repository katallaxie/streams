package streams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	assert.NotNil(t, r)
}

func TestNewMetric(t *testing.T) {
	m := NewMetrics()
	assert.NotNil(t, m)
}

func TestNewMonitor(t *testing.T) {
	m := NewMetrics()
	assert.NotNil(t, m)
	mon := NewMonitor(m)
	assert.NotNil(t, mon)
}
