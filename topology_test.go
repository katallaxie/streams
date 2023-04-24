package streams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNode(t *testing.T) {
	n := NewNode("test")
	assert.NotNil(t, n)
	assert.Equal(t, "test", n.Name())
	assert.Implements(t, (*Node)(nil), n)
}

func TestAddChild(t *testing.T) {
	n := NewNode("test")
	n.AddChild(NewNode("child"))
	assert.Equal(t, 1, len(n.Children()))

	nn := n.Children()
	assert.Equal(t, "child", nn[0].Name())
}

func TestNewTopology(t *testing.T) {
	top := NewTopology(NewNode("root"))
	assert.NotNil(t, top)
	assert.NotNil(t, top.Root())
	assert.Implements(t, (*Topology)(nil), top)

	assert.Equal(t, "root", top.Root().Name())
}
