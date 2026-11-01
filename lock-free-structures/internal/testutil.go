package internal

import "testing"

func CheckExactlyOnce(t *testing.T, dequeued [][]int, total int) {
	t.Helper()

	seen := make([]int, total)
	for _, values := range dequeued {
		for _, v := range values {
			seen[v]++
		}
	}

	for v, n := range seen {
		if n != 1 {
			t.Errorf("value %d dequeued %d times, want 1", v, n)
		}
	}
}
