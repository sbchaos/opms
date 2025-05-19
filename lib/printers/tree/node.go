package tree

import (
	"bytes"
	"fmt"
)

type EdgeType string

var (
	EdgeTypeEmpty EdgeType = " "
	EdgeTypeLink  EdgeType = "│"
	EdgeTypeMid   EdgeType = "├─"
	EdgeTypeEnd   EdgeType = "└─"
)

var indent = "    "

type Node[V any] struct {
	level int

	children map[string]*Node[V]
	Value    V
	Key      string
}

func (n *Node[V]) Level() int {
	return n.level
}

func NewNode[V any](name string, v V) *Node[V] {
	return &Node[V]{
		Value:    v,
		Key:      name,
		children: make(map[string]*Node[V]),
	}
}

func (n *Node[V]) AddNode(name string, n1 *Node[V]) {
	n1.level = n.level + 1
	n.children[name] = n1
}

func (n *Node[V]) AddChild(name string, v V) {
	n.children[name] = &Node[V]{
		Value:    v,
		children: make(map[string]*Node[V]),
		Key:      name,
		level:    n.level + 1,
	}
}

func (n *Node[V]) Bytes(buf *bytes.Buffer, parenPrefix string) []byte {
	num := len(n.children)
	i := 0
	for key, child := range n.children {
		i++
		curr := EdgeTypeMid
		paren := EdgeTypeLink
		childEnd := i == num
		if childEnd {
			curr = EdgeTypeEnd
			paren = EdgeTypeEmpty
		}
		fmt.Fprintf(buf, "%s%s %s [%v]\n", parenPrefix, curr, key, child.Value)

		prefix := fmt.Sprintf("%s%s%s", parenPrefix, paren, indent)
		child.Bytes(buf, prefix)
	}

	return buf.Bytes()
}

func (n *Node[V]) String() string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%s [%v]\n", n.Key, n.Value)
	return string(n.Bytes(buf, ""))
}
