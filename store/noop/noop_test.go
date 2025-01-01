package noop

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	n := New()
	require.NoError(t, n.Open())
}

func TestClose(t *testing.T) {
	n := New()
	require.NoError(t, n.Close())
}

func TestHas(t *testing.T) {
	n := New()
	ok, err := n.Has("foo")
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestGet(t *testing.T) {
	n := New()
	v, err := n.Get("foo")
	require.NoError(t, err)
	assert.Nil(t, v)
}

func TestSet(t *testing.T) {
	n := New()
	require.NoError(t, n.Set("foo", []byte("bar")))
}

func TestDelete(t *testing.T) {
	n := New()
	require.NoError(t, n.Delete("foo"))
}
