package resultext

import (
	"errors"
	"io"
	"testing"

	. "github.com/pchchv/go-assert"
)

type Struct struct{}

func TestResult(t *testing.T) {
	result := returnOk()
	Equal(t, true, result.IsOk())
	Equal(t, false, result.IsErr())
	Equal(t, true, result.Err() == nil)
	Equal(t, Struct{}, result.Unwrap())

	result = returnErr()
	Equal(t, false, result.IsOk())
	Equal(t, true, result.IsErr())
	Equal(t, false, result.Err() == nil)
	PanicMatches(t, func() {
		result.Unwrap()
	}, "Result.Unwrap(): result is Err")
}

func TestUnwrap(t *testing.T) {
	er := Err[int, error](io.EOF)
	PanicMatches(t, func() { er.Unwrap() }, "Result.Unwrap(): result is Err")

	v := er.UnwrapOr(3)
	Equal(t, 3, v)

	v = er.UnwrapOrElse(func() int { return 2 })
	Equal(t, 2, v)

	v = er.UnwrapOrDefault()
	Equal(t, 0, v)
}

func returnOk() Result[Struct, error] {
	return Ok[Struct, error](Struct{})
}

func returnErr() Result[Struct, error] {
	return Err[Struct, error](errors.New("bad"))
}
