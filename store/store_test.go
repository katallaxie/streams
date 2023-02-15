package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnimplemented(t *testing.T) {
	s := new(Unimplemented)

	assert.Error(t, ErrUnimplemented, s.Open())
	assert.Error(t, ErrUnimplemented, s.Close())

	ok, err := s.Has("foo")
	assert.Error(t, ErrUnimplemented, err)
	assert.False(t, ok)

	b, err := s.Get("foo")
	assert.Error(t, ErrUnimplemented, err)
	assert.Nil(t, b)

	assert.Equal(t, ErrUnimplemented, s.Set("foo", []byte("bar")))
	assert.Equal(t, ErrUnimplemented, s.Delete("foo"))
}
