package resultext

// Result represents the result of an operation that is successful or not.
type Result[T, E any] struct {
	ok   T
	err  E
	isOk bool
}

// Unwrap returns the values of the result.
// It panics if there is no result due to not checking for errors.
func (r Result[T, E]) Unwrap() T {
	if r.isOk {
		return r.ok
	}

	panic("Result.Unwrap(): result is Err")
}

// UnwrapOr returns the contained Ok value or a provided default.
func (r Result[T, E]) UnwrapOr(value T) T {
	if r.isOk {
		return r.ok
	}

	return value
}
