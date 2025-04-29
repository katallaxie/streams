package streams

import "github.com/katallaxie/pkg/slices"

var _ Node = (*node)(nil)

// Node is a node in a topology.
type Node interface {
	// AddChildren adds children to a node.
	AddChild(nodes ...Node)
	// Children returns the children of a node.
	Children() []Node
	// Name returns the name of a node.
	Name(names ...string) string
}

type node struct {
	id       int64
	name     string
	children []Node
}

// NewNode is a constructor for a new node in the topology.
func NewNode(name string) Node {
	n := new(node)
	n.name = name

	return n
}

// ID return the ID of a node.
func (n *node) ID() int64 {
	return n.id
}

// AddChild adds a child to a node.
func (n *node) AddChild(nodes ...Node) {
	n.children = append(n.children, nodes...)
}

// Children returns the children of a node.
func (n *node) Children() []Node {
	return n.children
}

// Name returns the name of a node.
func (n *node) Name(name ...string) string {
	if slices.GreaterThen(0, name...) {
		n.name = slices.First(name...)
	}

	return n.name
}

// Topology is a graph of nodes.
type Topology interface {
	// Root returns the root node of a topology.
	Root() Node
}

type topology struct {
	root Node
}

// NewTopology is a constructor for Topology.
func NewTopology(root Node) Topology {
	t := new(topology)
	t.root = root

	return t
}

// Root returns the root node of a topology.
func (t *topology) Root() Node {
	return t.root
}
