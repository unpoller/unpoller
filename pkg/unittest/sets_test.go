package unittest_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unpoller/unpoller/pkg/unittest"
)

func TestSets(t *testing.T) {
	s1 := unittest.NewSetFromSlice[string]([]string{"a", "b", "c", "c"})

	assert.Len(t, s1.Slice(), 3)
	assert.Contains(t, s1.Slice(), "a")
	assert.Contains(t, s1.Slice(), "b")
	assert.Contains(t, s1.Slice(), "c")

	s2 := unittest.NewSetFromMap[string](map[string]bool{
		"c": true,
		"d": false,
		"e": true,
	})

	assert.Len(t, s2.Slice(), 3)
	assert.Contains(t, s2.Slice(), "c")
	assert.Contains(t, s2.Slice(), "d")
	assert.Contains(t, s2.Slice(), "e")

	additions, deletions := s1.Difference(s2)

	assert.Len(t, additions, 2)
	assert.Len(t, deletions, 2)

	assert.Contains(t, additions, "a")
	assert.Contains(t, additions, "b")

	assert.Contains(t, deletions, "d")
	assert.Contains(t, deletions, "e")
}
