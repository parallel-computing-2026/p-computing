package queue

import "sync/atomic"

// node is a single node in LockFreeQueue.
// It holds a value and a pointer to the next node in the queue.
type node[T any] struct {
	value T
	next  atomic.Pointer[node[T]]
}

// LockFreeQueue is lock-free MPMC queue based on a linked list.
// Create with NewLockFreeQueue.
//
// It is safe for concurrent use.
type LockFreeQueue[T any] struct {
	head atomic.Pointer[node[T]]
	tail atomic.Pointer[node[T]]
}

// NewLockFreeQueue initializes a new LockFreeQueue of type T.
//
// It is safe for concurrent use.
func NewLockFreeQueue[T any]() *LockFreeQueue[T] {
	dummy := &node[T]{}

	queue := &LockFreeQueue[T]{}
	queue.head.Store(dummy)
	queue.tail.Store(dummy)

	return queue
}

// Enqueue adds a new element to the tail of the queue.
func (l *LockFreeQueue[T]) Enqueue(value T) {
	n := &node[T]{value: value}
	var tail *node[T]
	for {
		tail = l.tail.Load()
		next := tail.next.Load()

		if tail == l.tail.Load() {
			// no one enqueued next, so we try to do it
			if next == nil {
				if tail.next.CompareAndSwap(next, n) {
					break
				}
			} else {
				// someone else enqueued next, so we help them
				l.tail.CompareAndSwap(tail, next)
			}
		}
	}

	// same as else branch above
	l.tail.CompareAndSwap(tail, n)
}

// Dequeue removes the element at the head of the queue and returns its value.
// It returns false if the queue is empty.
func (l *LockFreeQueue[T]) Dequeue() (T, bool) {
	var zero T
	for {
		head := l.head.Load()
		tail := l.tail.Load()
		next := head.next.Load()

		if head != l.head.Load() {
			// out of sync with the current state
			continue
		}

		if tail == head {
			if next == nil {
				// queue is empty
				return zero, false
			}
			// someone else enqueued next, so we help them
			l.tail.CompareAndSwap(tail, next)
		} else {
			if next == nil {
				// someone else dequeued next, so we try again
				continue
			}

			if l.head.CompareAndSwap(head, next) {
				return next.value, true
			}
		}
	}
}
