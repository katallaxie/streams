package streams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStream(t *testing.T) {
	src := newMockSource[string, string]()
	s := NewStream[string, string](src)

	assert.NotNil(t, s)
}

func TestCountMetric(t *testing.T) {
	m := newCountMetric("foo")
	m.inc()

	assert.Equal(t, m.value, float64(1))
	assert.Equal(t, m.nodeName, "foo")

	ch := make(chan Metric)

	go func() {
		m.Collect(ch)
	}()

	mm := <-ch
	assert.Equal(t, mm, m)
	assert.Equal(t, m.value, float64(1))
}
