package stack

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestConcurrentWrites(t *testing.T) {
	stack := NewLockFreeStack[int64]()

	const writers = 4
	const numbers = 10000

	var counter atomic.Int64
	var wg sync.WaitGroup
	for range writers {
		wg.Go(func() {
			for range numbers {
				stack.Push(counter.Add(1) - 1)
			}
		})
	}
	wg.Wait()

	values := make(map[int64]struct{}, numbers*writers)
	for v, ok := stack.Pop(); ok; v, ok = stack.Pop() {
		if _, ok := values[v]; ok {
			t.Errorf("got duplicate value: %d", v)
		}
		values[v] = struct{}{}
	}

	if len(values) != numbers*writers {
		t.Fatalf("got %d values, expected %d", len(values), numbers)
	}
}

func TestConcurrentRW(t *testing.T) {
	stack := NewLockFreeStack[int64]()

	const writers = 4
	const readers = 3
	const numbers = 10000

	var writerCounter atomic.Int64
	var wg sync.WaitGroup
	for range writers {
		wg.Go(func() {
			for range numbers {
				stack.Push(writerCounter.Add(1) - 1)
			}
		})
	}

	var readerCounter atomic.Int64
	for range readers {
		wg.Go(func() {
			// can potentially be stuck in an endless loop
			// prolly needs a timeout of some sort
			for readerCounter.Load() < numbers*writers {
				_, ok := stack.Pop()
				if ok {
					readerCounter.Add(1)
				}
			}
		})
	}
	wg.Wait()

	if readerCounter != writerCounter {
		t.Errorf("got %d, want %d", readerCounter.Load(), writerCounter.Load())
	}
}
