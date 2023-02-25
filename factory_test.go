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
