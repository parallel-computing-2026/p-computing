package stack

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/parallel-computing-2026/p-computing/lock-free-structures/internal"
)

const (
	writers = 4
	readers = 3
	numbers = 10000
	total   = writers * numbers
)

func push(stack *LockFreeStack[int], wg *sync.WaitGroup) {
	for w := range writers {
		wg.Go(func() {
			for i := range numbers {
				stack.Push(w*numbers + i)
			}
		})
	}
}

func TestConcurrentWrites(t *testing.T) {
	stack := NewLockFreeStack[int]()

	var wg sync.WaitGroup
	push(stack, &wg)
	wg.Wait()

	values := make([]int, 0, total)
	for val, ok := stack.Pop(); ok; val, ok = stack.Pop() {
		values = append(values, val)
	}

	internal.CheckExactlyOnce(t, [][]int{values}, total)
}

func TestConcurrentRW(t *testing.T) {
	stack := NewLockFreeStack[int]()

	var wg sync.WaitGroup
	push(stack, &wg)

	var count atomic.Int64
	values := make([][]int, readers)
	for r := range readers {
		wg.Go(func() {
			for count.Load() < total {
				if val, ok := stack.Pop(); ok {
					count.Add(1)
					values[r] = append(values[r], val)
				}
			}
		})
	}
	wg.Wait()

	internal.CheckExactlyOnce(t, values, total)
}
