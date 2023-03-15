package leveldb

import (
	"testing"

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
}

func TestHas(t *testing.T) {
	l := New()
	assert.NotNil(t, l)

	err := l.Open()
	assert.NoError(t, err)

	has, err := l.Has("foo")
	assert.Error(t, err)
	assert.False(t, has)
}

func TestGet(t *testing.T) {
	l := New()
	assert.NotNil(t, l)

	err := l.Open()
	assert.NoError(t, err)

	_, err = l.Get("foo")
	assert.Error(t, err)
}

func TestSet(t *testing.T) {
	l := New()
	assert.NotNil(t, l)

	err := l.Open()
	assert.NoError(t, err)

	err = l.Set("foo", []byte("bar"))
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	l := New()
	assert.NotNil(t, l)

	err := l.Open()
	assert.NoError(t, err)

	err = l.Set("foo", []byte("bar"))
	assert.NoError(t, err)

	err = l.Delete("foo")
	assert.NoError(t, err)
}
