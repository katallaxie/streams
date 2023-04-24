package streams

var _ Node = (*node)(nil)

type node struct {
	name     string
	children []Node
}

// NewNode is a constructor for a new node in the topology.
func NewNode(name string) Node {
	n := new(node)
	n.name = name

	return n
}

// AddChild adds a child to a node.
func (n *node) AddChild(child Node) {
	n.children = append(n.children, child)
}

// Children returns the children of a node.
func (n *node) Children() []Node {
	return n.children
}

// Name returns the name of a node.
func (n *node) Name() string {
	return n.name
}

// Topology is a graph of nodes.
type Topology interface {
	// Root returns the root node of a topology.
	Root() Node
}

// Node is a node in a topology.
type Node interface {
	// AddChild adds a child to a node.
	AddChild(Node)

	// Children returns the children of a node.
	Children() []Node

	// Name returns the name of a node.
	Name() string
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
