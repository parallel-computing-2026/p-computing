package stack

import "sync/atomic"

// node is a single node in LockFreeStack.
// It holds a value and a pointer to the next node in the stack.
type node[T any] struct {
	value T
	next  *node[T]
}

// LockFreeStack is lock-free treiber stack based on a linked list.
// Create with NewLockFreeStack.
//
// It is safe for concurrent use without.
type LockFreeStack[T any] struct {
	head atomic.Pointer[node[T]] // defaults to nil
}

// NewLockFreeStack initializes a new LockFreeStack of type T.
//
// It is safe for concurrent use without.
func NewLockFreeStack[T any]() *LockFreeStack[T] {
	return &LockFreeStack[T]{}
}

// Push inserts a new element at the top of the stack.
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

// Pop removes the element at the top of the stack and returns its value.
// It returns false if the stack is empty.
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
