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
