package trie_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/sbchaos/opms/lib/trie"
)

func TestRuneTrieContains(t *testing.T) {
	trie := trie.NewTrie[string]()

	trie.Insert("c", "c")
	trie.Insert("apple", "apple")
	trie.Insert("banana", "banana")

	cases := []struct {
		word   string
		exists bool
	}{
		{"c", true},
		{"cc", false},
		{"ce", false},
		{"banana", true},
		{"app", false},
		{"aple", false},
		{"ban", false},
		{"apple", true},
		{"apple1", false},
		{"aaapple", false},
	}

	for _, c := range cases {
		actual := trie.Contains(c.word)
		if actual != c.exists {
			if c.exists {
				t.Errorf("%s is expected to be found in the trie", c.word)
			} else {
				t.Errorf("%s is not expected to be found in the trie", c.word)
			}
		}
	}
}

func TestRuneTrieStartsWith(t *testing.T) {
	trie := trie.NewTrie[string]()

	trie.Insert("apple", "apple")
	trie.Insert("aplhabet", "aplhabet")
	trie.Insert("tree", "tree")

	if !trie.StartsWith("a") {
		t.Error("trie StartsWith(a) must return true, but returned false")
	}

	if trie.StartsWith("try") {
		t.Error("trie StartsWith(try) must return false, but returned true")
	}
}

func TestRuneTrieSearchByPrefix(t *testing.T) {
	trie := trie.NewTrie[string]()

	trie.Insert("c", "c")
	trie.Insert("apple", "apple")
	trie.Insert("banana", "banana")
	trie.Insert("alphabet", "alphabet")
	trie.Insert("alcohol", "alcohol")

	actual := trie.Search("a")
	sort.Strings(actual)

	expected := []string{"alcohol", "alphabet", "apple"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%v != %v", actual, expected)
	}
}
