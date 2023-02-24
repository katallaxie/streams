package noop

import (
	"testing"

	"github.com/ionos-cloud/streams/msg"
	"github.com/stretchr/testify/assert"
)

func TestNewSink(t *testing.T) {
	t.Parallel()

	s := NewSink[string, string]()
	assert.NotNil(t, s)
}

func TestSinkWrite(t *testing.T) {
	t.Parallel()

	s := NewSink[string, string]()
	assert.NotNil(t, s)

	err := s.Write(msg.NewMessage("foo", "bar", 0, 0, "", nil))
	assert.Nil(t, err)
}
