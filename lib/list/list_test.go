package list_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sbchaos/opms/lib/list"
)

func TestRemove(t *testing.T) {
	t.Run("remove", func(t *testing.T) {
		lst := []string{"a", "b", "c"}
		lst = list.Remove(lst, "b")

		assert.Equal(t, []string{"a", "c"}, lst)
	})
}
