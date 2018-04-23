package someprocess

import (
	"testing"
)

func TestRun(t *testing.T) {
	call = func(val int) int {
		return val + 10
	}

	patterns := []struct {
		val      int
		expected int
	}{
		{2, 12},
		{8, 18},
		{-10, 0},
	}

	for idx, pattern := range patterns {
		actual := Run(pattern.val)
		if pattern.expected != actual {
			t.Errorf("pattern %d: want %d, actual %d", idx, pattern.expected, actual)
		}
	}
}
