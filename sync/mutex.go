package syncext

import "sync"

// Mutex creates a type safe mutex wrapper ensuring one cannot access the
// values of a locked values without first gaining a lock.
type Mutex[T any] struct {
	m     *sync.Mutex
	value T
}
