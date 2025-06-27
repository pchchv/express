package syncext

import "sync"

// MutexGuard protects the inner contents of a Mutex for safety and unlocking.
type MutexGuard[T any, M interface{ Unlock() }] struct {
	m M
	T T // is the inner generic type of the Mutex
}

// Unlock unlocks the Mutex value.
func (g MutexGuard[T, M]) Unlock() {
	g.m.Unlock()
}

// Mutex creates a type safe mutex wrapper ensuring one cannot access the
// values of a locked values without first gaining a lock.
type Mutex[T any] struct {
	m     *sync.Mutex
	value T
}

// NewMutex creates a new Mutex for use.
func NewMutex[T any](value T) Mutex[T] {
	return Mutex[T]{
		m:     new(sync.Mutex),
		value: value,
	}
}

// Lock locks the Mutex and returns value for use,
// safe for mutation if the lock is already in use,
// the calling goroutine blocks until the mutex is available.
func (m Mutex[T]) Lock() MutexGuard[T, *sync.Mutex] {
	m.m.Lock()
	return MutexGuard[T, *sync.Mutex]{
		m: m.m,
		T: m.value,
	}
}

// Unlock unlocks the Mutex accepting a value to set as the new or mutated value.
// It is optional to pass a new value to be set but NOT required for there reasons:
// 1. If the internal value is already mutable then no need to set as changes apply as they happen.
// 2. If there's a failure working with the locked value you may NOT want to set it, but still unlock.
// 3. Supports locked values that are not mutable.
//
// It is a run-time error if the Mutex is not locked on entry to Unlock.
func (m Mutex[T]) Unlock() {
	m.m.Unlock()
}
