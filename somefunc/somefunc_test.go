package somefunc

import (
	"testing"
)

func TestRun(t *testing.T) {

	patterns := []struct {
		val      int
		expected int
	}{
		{2, 12},
		{8, 18},
		{-10, 0},
	}

	for idx, pattern := range patterns {
		// Clientのnewの際に、モックオブジェクトを引数にする
		c := Client{&mockCaller{}}
		actual := c.Run(pattern.val)
		if pattern.expected != actual {
			t.Errorf("pattern %d: want %d, actual %d", idx, pattern.expected, actual)
		}
	}
}

type mockCaller struct{}

// モックの方のメソッドでは引数+10を返却する
func (s *mockCaller) call(val int) int {
	return val + 10
}
