package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnimplemented(t *testing.T) {
	s := new(Unimplemented)

	require.ErrorIs(t, ErrUnimplemented, s.Open())
	require.ErrorIs(t, ErrUnimplemented, s.Close())

	ok, err := s.Has("foo")
	require.ErrorIs(t, ErrUnimplemented, err)
	assert.False(t, ok)

	b, err := s.Get("foo")
	require.ErrorIs(t, ErrUnimplemented, err)
	assert.Nil(t, b)

	require.ErrorIs(t, ErrUnimplemented, s.Set("foo", []byte("bar")))
	require.ErrorIs(t, ErrUnimplemented, s.Delete("foo"))
}
