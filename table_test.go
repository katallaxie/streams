package streams

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnimplemented_Set(t *testing.T) {
	tbl := &tableUnimplemented{}

	err := tbl.Set("key", []byte("value"))
	require.Error(t, err)
	assert.Equal(t, ErrNotImplemented, err)
}

func TestUnimplemented_Delete(t *testing.T) {
	tbl := &tableUnimplemented{}

	err := tbl.Delete("key")
	require.Error(t, err)
	assert.Equal(t, ErrNotImplemented, err)
}

func TestUnimplemented_Setup(t *testing.T) {
	tbl := &tableUnimplemented{}

	err := tbl.Setup()
	require.Error(t, err)
	assert.Equal(t, ErrNotImplemented, err)
}

func TestUnimplemented_Error(t *testing.T) {
	tbl := &tableUnimplemented{}

	err := tbl.Error()
	require.Error(t, err)
	assert.Equal(t, ErrNotImplemented, err)
}
