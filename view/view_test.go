package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	v := New[string](nil, nil, nil, nil)
	assert.NotNil(t, v)
}
