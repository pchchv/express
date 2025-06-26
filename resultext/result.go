package resultext

// Result represents the result of an operation that is successful or not.
type Result[T, E any] struct {
	ok   T
	err  E
	isOk bool
}
