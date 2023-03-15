package leveldb

import (
	"testing"

	"github.com/ionos-cloud/streams/codec"
	"github.com/ionos-cloud/streams/msg"
	"github.com/ionos-cloud/streams/store"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	l := New()
	assert.NotNil(t, l)
}

func TestOpen(t *testing.T) {
	l := New()
	assert.NotNil(t, l)

	err := l.Open()
	assert.NoError(t, err)
	defer l.Close()
}

func TestHas(t *testing.T) {
	l := New()
	assert.NotNil(t, l)

	err := l.Open()
	assert.NoError(t, err)
	defer l.Close()

	has, err := l.Has("foo")
	assert.Error(t, err)
	assert.False(t, has)
}

func TestGet(t *testing.T) {
	l := New()
	assert.NotNil(t, l)

	err := l.Open()
	assert.NoError(t, err)
	defer l.Close()

	_, err = l.Get("foo")
	assert.Error(t, err)
}

func TestSet(t *testing.T) {
	l := New()
	assert.NotNil(t, l)

	err := l.Open()
	assert.NoError(t, err)
	defer l.Close()

	err = l.Set("foo", []byte("bar"))
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	l := New()
	assert.NotNil(t, l)

	err := l.Open()
	assert.NoError(t, err)
	defer l.Close()

	err = l.Set("foo", []byte("bar"))
	assert.NoError(t, err)

	err = l.Delete("foo")
	assert.NoError(t, err)
}

func TestSink(t *testing.T) {
	l := New()
	assert.NotNil(t, l)

	err := l.Open()
	assert.NoError(t, err)
	defer l.Close()

	s := store.NewSink(l, codec.StringEncoder)
	assert.NotNil(t, s)

	err = s.Write(msg.NewMessage("foo", "bar", 0, 0, "bar", nil))
	assert.NoError(t, err)

	b, err := l.Get("foo")
	assert.NoError(t, err)
	assert.Equal(t, "bar", string(b))
}
