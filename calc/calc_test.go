package calc

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("before test")
	code := m.Run()
	fmt.Println("after test")
	os.Exit(code)
}

func TestAdd(t *testing.T) {

	patterns := []struct {
		a        int
		b        int
		expected int
	}{
		{1, 2, 3},
		{10, -2, 8},
		{-10, -2, -12},
	}

	for idx, pattern := range patterns {
		actual := Add(pattern.a, pattern.b)
		if pattern.expected != actual {
			t.Errorf("pattern %d: want %d, actual %d", idx, pattern.expected, actual)
		}
	}
}
