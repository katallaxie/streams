package table

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTopic(t *testing.T) {
	topic := Topic("foo")
	assert.Equal(t, "table.foo", topic.String())
}
