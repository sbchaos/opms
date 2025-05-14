package trie

type runeNode[V any] struct {
	children map[rune]*runeNode[V]
	last     bool
	Value    V // Only the last node has value
}

type Trie[V any] struct {
	root *runeNode[V]
	size int
}

func NewTrie[V any]() *Trie[V] {
	return &Trie[V]{root: &runeNode[V]{children: make(map[rune]*runeNode[V]), last: false}}
}

func (t *Trie[V]) Insert(word string, value V) bool {
	exists := true
	current := t.root
	for _, letter := range word {
		n, ok := current.children[letter]
		if !ok {
			exists = false

			n = &runeNode[V]{children: make(map[rune]*runeNode[V]), last: false}
			current.children[letter] = n
		}
		current = n
	}
	current.last = true
	current.Value = value

	if !exists {
		t.size++
	}

	return exists
}

func (t *Trie[V]) Search(word string) []string {
	node, r := t.nodeByPrefix(word)

	return search(node, r, []rune(word[:len(word)-1]))
}

func (t *Trie[V]) StartsWith(prefix string) bool {
	node, _ := t.nodeByPrefix(prefix)

	return node != nil
}

func (t *Trie[V]) Contains(word string) bool {
	n, _ := t.nodeByPrefix(word)

	return n != nil && n.last
}

func (t *Trie[V]) nodeByPrefix(prefix string) (*runeNode[V], rune) {
	current := t.root
	var r rune
	for _, letter := range prefix {
		n, ok := current.children[letter]
		if !ok {
			return nil, 0
		}

		current = n
		r = letter
	}

	return current, r
}

func search[V any](currentNode *runeNode[V], r rune, prefix []rune) []string {
	words := []string{}
	if currentNode == nil {
		return words
	}

	newPrefix := append(prefix, r)
	if currentNode.last {
		words = append(words, string(newPrefix))
	}

	for letter, node := range currentNode.children {
		newWords := search(node, letter, newPrefix)
		words = append(words, newWords...)
	}

	return words
}
