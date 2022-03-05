package rope

import (
	"testing"
	"math"
	"fmt"
)

func assertSameValue[T comparable](t *testing.T, a, b *Rope[T]) {
	aValue := a.Value()
	bValue := b.Value()
	if len(aValue) != len(bValue) {
		t.Errorf("Expected equal length ropes\n1: %v\n2: %v\n", aValue, bValue)
		return
	}
	for i := range aValue {
		if aValue[i] != bValue[i] {
			t.Errorf("Expected equal value ropes\n1: %v\n2: %v\n", aValue, bValue)
			return
		}
	}
}

func assertValue[T comparable](t *testing.T, rope *Rope[T], value[]T) {
	ropeValue := rope.Value()
	if len(ropeValue) != len(value) {
		t.Errorf("Different length\n1: %v\n2: %v\n", ropeValue, value)
		return
	}
	for i := range value {
		if ropeValue[i] != value[i] {
			t.Errorf("Expected equal value ropes\n1: %v\n2: %v\n", ropeValue, value)
			return
		}
	}
}

func assert(t *testing.T, value bool, log... interface{}) {
	if !value {
		t.Error(log...)
	}
}

func maxDepth[T any](rope *Rope[T]) int {
	if rope.value != nil {
		return 1
	}
	leftDepth := maxDepth(rope.left)
	rightDepth := maxDepth(rope.right)
	if leftDepth > rightDepth {
		return 1 + leftDepth
	}
	return 1 + rightDepth
}

var testSettings = &Settings {
	SplitLength: 4,
	JoinLength:  2,
	Rebalance:   1.001,
}

func TestInsert(t *testing.T) {
	originalValue := []int{0, 1, 2, 3, 4, 5, 6, 7}
	rope := NewRope(originalValue, testSettings)
	newRope := rope.Insert(2, []int{-1, -2, -3})

	assertValue[int](t, rope, originalValue)
	assertValue[int](t, newRope, []int {
		0, 1, -1, -2, -3, 2, 3, 4, 5, 6, 7,
	})
}

func TestDelete(t *testing.T) {
	originalValue := []int{0, 1, 2, 3, 4, 5, 6, 7}
	rope := NewRope(originalValue, testSettings)
	newRope := rope.Remove(1, 6)

	assertValue[int](t, rope, originalValue)
	assertValue[int](t, newRope, []int {
		0, 6, 7,
	})
}

func TestSlice(t *testing.T) {
	originalValue := []int{0, 1, 2, 3, 4, 5, 6, 7}
	rope := NewRope(originalValue, testSettings)
	newRope := rope.Remove(1, 6)

	assertValue[int](t, rope, originalValue)
	assertValue[int](t, newRope, []int {
		0, 6, 7,
	})
}

func TestRebalance(t *testing.T) {
	const n = 1000
	originalValue := []int{}
	rope := NewRope(originalValue, testSettings)
	newRope := rope

	for i := 0; i < n; i ++ {
		newRope = newRope.Insert(0, []int{0, 1, 2, 3, 4, 5, 6, 7})
	}

	assert(t, maxDepth[int](newRope) >= 10, "Max depth too low:", maxDepth[int](newRope))

	balancedRope := NewRope(newRope.Value(), testSettings)
	balancedRope.Rebalance()

	assert(t, maxDepth[int](balancedRope) <= int(math.Log2(n*8)), "Rebalance didn't fix max depth:", maxDepth[int](balancedRope))

	assertValue[int](t, rope, originalValue)
	assertSameValue[int](t, balancedRope, newRope)
}

var inputs = []int{1, 10, 100, 1000, 10000, 100000}

func BenchmarkRopeInsert(b *testing.B) {
	for _, input := range inputs {
		b.Run(fmt.Sprintf("rope_%v_insertions", input), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rope := NewRope([]byte{'a', 'b', 'c', 'd'}, DefaultSettings)
				for j := 0; j < input;j ++ {
					index := (j * 77777777) % rope.Length()
					if index < 0 {
						index = -index
					}
					rope = rope.Insert(index, []byte{'a', 'b', 'c', 'd'})
				}
			}
		})
	}
}

func BenchmarkRopeInsertRebalance(b *testing.B) {
	for _, input := range inputs {
		b.Run(fmt.Sprintf("rope_%v_insertions", input), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rope := NewRope([]byte{'a', 'b', 'c', 'd'}, DefaultSettings)
				for j := 0; j < input;j ++ {
					if j % 1000 == 0 {
						rope.Rebalance()
					}
					index := (j * 77777777) % rope.Length()
					if index < 0 {
						index = -index
					}
					rope = rope.Insert(index, []byte{'a', 'b', 'c', 'd'})
				}
			}
		})
	}
}

func BenchmarkStringInsert(b *testing.B) {
	for _, input := range inputs {
		b.Run(fmt.Sprintf("rope_%v_insertions", input), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				str := "abcd"
				for j := 0; j < input;j ++ {
					index := (j * 77777777) % len(str)
					if index < 0 {
						index = -index
					}
					str = str[:index] + "abcd" + str[index:]
				}
			}
		})
	}
}
