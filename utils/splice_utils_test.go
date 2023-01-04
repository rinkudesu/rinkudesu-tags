package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContains_ValueNotPresent_ReturnsFalse(t *testing.T) {
	haystack := []string{"aaa", "bbb", "ccc"}
	needle := "ddd"

	result := Contains(haystack, needle)

	assert.False(t, result)
}

func TestContains_ValuePresent_ReturnsTrue(t *testing.T) {
	haystack := []string{"aaa", "bbb", "ccc"}
	needles := []string{"aaa", "bbb", "ccc"}

	for _, needle := range needles {
		t.Run(needle, func(t *testing.T) {
			t.Parallel()
			assert.True(t, Contains(haystack, needle))
		})
	}

}
