package optionext

// Option represents a values that represents a values existence.
//
// nil is usually used on Go however this has two problems:
// 1. Checking if the return values is nil is NOT enforced and can lead to panics.
// 2. Using nil is not good enough when nil itself is a valid value.
//
// This implements the sql.Scanner interface and can be used as a sql value for reading and writing.
// It supports:
// - String
// - Bool
// - Uint8
// - Float64
// - Int16
// - Int32
// - Int64
// - interface{}/any
// - time.Time
// - Struct - when type is convertable to []byte and assumes JSON.
// - Slice - when type is convertable to []byte and assumes JSON.
// - Map types - when type is convertable to []byte and assumes JSON.
//
// This also implements the `json.Marshaler` and `json.Unmarshaler` interfaces.
// The only caveat is a None value will result in a JSON `null` value.
type Option[T any] struct {
	value  T
	isSome bool
}
