package noop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSource(t *testing.T) {
	s := NewSource[string, string](nil)
	assert.NotNil(t, s)
}

func TestCommit(t *testing.T) {
	s := NewSource[string, string](nil)
	assert.NotNil(t, s)

	err := s.Commit()
	assert.Nil(t, err)
}

func TestClose(t *testing.T) {
	s := NewSource[string, string](nil)
	assert.NotNil(t, s)

	s.Close()
}

func TestMessages(t *testing.T) {
	s := NewSource[string, string](nil)
	assert.NotNil(t, s)

	out := s.Messages()
	assert.NotNil(t, out)
}
