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

var (
	bold  = "\033[1;"
	reset = "\033[0m"

	indent = "    "
)

type Node[V any] struct {
	level int

	children []*Node[V]
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
		children: []*Node[V]{},
	}
}

func (n *Node[V]) AddNode(n1 *Node[V]) {
	n1.level = n.level + 1
	n.children = append(n.children, n1)
}

func (n *Node[V]) AddChild(name string, v V) {
	n.children = append(n.children, &Node[V]{
		Value:    v,
		children: []*Node[V]{},
		Key:      name,
		level:    n.level + 1,
	})
}

func (n *Node[V]) Bytes(buf *bytes.Buffer, parenPrefix string) []byte {
	num := len(n.children)
	i := 0
	for _, child := range n.children {
		i++
		curr := EdgeTypeMid
		paren := EdgeTypeLink
		childEnd := i == num
		if childEnd {
			curr = EdgeTypeEnd
			paren = EdgeTypeEmpty
		}
		content := fmt.Sprintf("%s [%v]", child.Key, child.Value)
		fmt.Fprintf(buf, "%s%s %s \n", parenPrefix, curr, content)

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
