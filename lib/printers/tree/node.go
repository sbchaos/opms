package tree

import (
	"fmt"
	"io"
	"os"

	"github.com/sbchaos/opms/lib/color"
	"github.com/sbchaos/opms/lib/term"
)

type EdgeType string

var (
	EdgeTypeEmpty EdgeType = " "
	EdgeTypeLink  EdgeType = "│"
	EdgeTypeMid   EdgeType = "├─"
	EdgeTypeEnd   EdgeType = "└─"
)

var (
	indent = "    "
)

type Tree[V any] struct {
	root *Node[V]
	schm color.Scheme
}

func NewTreeWithAutoDetect[V any]() *Tree[V] {
	t := term.FromEnv(0, 0)
	scheme := color.NewColorScheme(t.IsColorEnabled(), t.Is256ColorSupported(), t.IsTrueColorSupported())
	return NewTree[V](scheme)
}

func NewTree[V any](scheme color.Scheme) *Tree[V] {
	var zero V
	return &Tree[V]{
		schm: scheme,
		root: NewNode("Root", zero),
	}
}

func (t *Tree[V]) Root() *Node[V] {
	return t.root
}

func (t *Tree[V]) Render(w io.Writer) {
	if len(t.root.children) == 0 {
		fmt.Fprintf(w, t.schm.Colorize(color.LightRed, "", "<Tree is Empty>\n"))
	}

	f := t.root.children[0]
	fmt.Fprintf(w, t.schm.Colorize(f.Color, f.Style, f.Content()))
	fmt.Fprintf(w, "\n")
	f.Print(w, t.schm, "")
}

type Node[V any] struct {
	level int

	children []*Node[V]
	Value    V
	Key      string

	Color int
	Style string
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

func (n *Node[V]) Print(w io.Writer, schm color.Scheme, parenPrefix string) {
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

		toShow := schm.Colorize(child.Color, child.Style, child.Content())
		fmt.Fprintf(w, "%s%s %s\n", parenPrefix, curr, toShow)

		prefix := fmt.Sprintf("%s%s%s", parenPrefix, paren, indent)
		child.Print(w, schm, prefix)
	}
}

func (n *Node[V]) String() {
	fmt.Fprintf(os.Stdout, "%s [%v]\n", n.Key, n.Value)
	noScheme := color.NewColorScheme(false, false, false)
	n.Print(os.Stdout, noScheme, "")
}

func (n *Node[V]) Content() string {
	return fmt.Sprintf("%s [%v]", n.Key, n.Value)
}
