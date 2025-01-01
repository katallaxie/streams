package leveldb

import (
	"testing"

	"github.com/katallaxie/streams/codec"
	"github.com/katallaxie/streams/msg"
	"github.com/katallaxie/streams/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	l := New()
	assert.NotNil(t, l)
}

func TestOpen(t *testing.T) {
	l := New()
	assert.NotNil(t, l)

	err := l.Open()
	require.NoError(t, err)
	defer l.Close()
}

func TestHas(t *testing.T) {
	l := New()
	assert.NotNil(t, l)

	err := l.Open()
	require.NoError(t, err)
	defer l.Close()

	has, err := l.Has("foo")
	require.Error(t, err)
	assert.False(t, has)
}

func TestGet(t *testing.T) {
	l := New()
	assert.NotNil(t, l)

	err := l.Open()
	require.NoError(t, err)
	defer l.Close()

	_, err = l.Get("foo")
	require.Error(t, err)
}

func TestSet(t *testing.T) {
	l := New()
	assert.NotNil(t, l)

	err := l.Open()
	require.NoError(t, err)
	defer l.Close()

	err = l.Set("foo", []byte("bar"))
	require.NoError(t, err)
}

func TestDelete(t *testing.T) {
	l := New()
	assert.NotNil(t, l)

	err := l.Open()
	require.NoError(t, err)
	defer l.Close()

	err = l.Set("foo", []byte("bar"))
	require.NoError(t, err)

	err = l.Delete("foo")
	require.NoError(t, err)
}

func TestSink(t *testing.T) {
	l := New()
	assert.NotNil(t, l)

	err := l.Open()
	require.NoError(t, err)
	defer l.Close()

	s := store.NewSink(l, codec.StringEncoder)
	assert.NotNil(t, s)

	err = s.Write(msg.NewMessage("foo", "bar", 0, 0, "bar", nil))
	require.NoError(t, err)

	b, err := l.Get("foo")
	require.NoError(t, err)
	assert.Equal(t, "bar", string(b))
}
