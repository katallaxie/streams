package streams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnimplemented_Set(t *testing.T) {
	tbl := &tableUnimplemented{}

	err := tbl.Set("key", []byte("value"))
	assert.Error(t, err)
	assert.Equal(t, ErrNotImplemented, err)
}

func TestUnimplemented_Delete(t *testing.T) {
	tbl := &tableUnimplemented{}

	err := tbl.Delete("key")
	assert.Error(t, err)
	assert.Equal(t, ErrNotImplemented, err)
}

func TestUnimplemented_Setup(t *testing.T) {
	tbl := &tableUnimplemented{}

	err := tbl.Setup()
	assert.Error(t, err)
	assert.Equal(t, ErrNotImplemented, err)
}

func TestUnimplemented_Error(t *testing.T) {
	tbl := &tableUnimplemented{}

	err := tbl.Error()
	assert.Error(t, err)
	assert.Equal(t, ErrNotImplemented, err)
}
