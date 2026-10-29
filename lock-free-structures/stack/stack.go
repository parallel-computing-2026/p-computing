package stack

import "sync/atomic"

type node[T any] struct {
	value T
	next  *node[T]
}

type LockFreeStack[T any] struct {
	head atomic.Pointer[node[T]] // defaults to nil
}

func NewLockFreeStack[T any]() *LockFreeStack[T] {
	return &LockFreeStack[T]{}
}

func (l *LockFreeStack[T]) Push(value T) {
	n := &node[T]{value: value}
	for {
		old := l.head.Load()
		n.next = old
		// check that the head did not change and swap to
		// the new node in a single operation
		if l.head.CompareAndSwap(old, n) {
			return
		}
	}
}

func (l *LockFreeStack[T]) Pop() (T, bool) {
	var zero T
	for {
		old := l.head.Load()
		// we can safely exit here
		// when the stack is empty
		if old == nil {
			return zero, false
		}

		next := old.next
		// swap head to the next element if it did not change
		if l.head.CompareAndSwap(old, next) {
			return old.value, true
		}
	}
}
