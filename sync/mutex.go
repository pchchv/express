package syncext

import (
	"sync"

	"github.com/pchchv/express/resultext"
)

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

// PerformMut safely locks and unlocks the Mutex values and performs the provided function returning its error if one
// otherwise setting the returned value as the new mutex value.
func (m Mutex[T]) PerformMut(f func(T)) {
	guard := m.Lock()
	f(guard.T)
	guard.Unlock()
}

// TryLock tries to lock Mutex and reports whether it succeeded.
// If it does the value is returned for use in the Ok result otherwise Err with empty value.
func (m Mutex[T]) TryLock() resultext.Result[MutexGuard[T, *sync.Mutex], struct{}] {
	if m.m.TryLock() {
		return resultext.Ok[MutexGuard[T, *sync.Mutex], struct{}](MutexGuard[T, *sync.Mutex]{
			m: m.m,
			T: m.value,
		})
	} else {
		return resultext.Err[MutexGuard[T, *sync.Mutex], struct{}](struct{}{})
	}
}

// RMutexGuard protects the inner contents of a RWMutex for safety and unlocking.
type RMutexGuard[T any] struct {
	rw *sync.RWMutex
	// T is the inner generic type of the Mutex
	T T
}
