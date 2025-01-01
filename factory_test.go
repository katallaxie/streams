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

	assert.InEpsilon(t, 1, m.value, 1)
	assert.Equal(t, "foo", m.nodeName)

	ch := make(chan Metric)

	go func() {
		m.Collect(ch)
	}()

	mm := <-ch
	assert.Equal(t, mm, m)
	assert.InEpsilon(t, 1, m.value, 1)
}
