package optionext

import (
	"encoding/json"
	"testing"
	"time"

	. "github.com/pchchv/go-assert"
)

type testStruct struct{}

func TestAndXXX(t *testing.T) {
	s := Some(1)
	Equal(t, Some(3), s.And(func(i int) int { return 3 }))
	Equal(t, Some(3), s.AndThen(func(i int) Option[int] { return Some(3) }))
	Equal(t, None[int](), s.AndThen(func(i int) Option[int] { return None[int]() }))

	n := None[int]()
	Equal(t, None[int](), n.And(func(i int) int { return 3 }))
	Equal(t, None[int](), n.AndThen(func(i int) Option[int] { return Some(3) }))
	Equal(t, None[int](), n.AndThen(func(i int) Option[int] { return None[int]() }))
	Equal(t, None[int](), s.AndThen(func(i int) Option[int] { return None[int]() }))
}

func TestUnwraps(t *testing.T) {
	none := None[int]()
	PanicMatches(t, func() { none.Unwrap() }, "Option.Unwrap: option is None")

	v := none.UnwrapOr(3)
	Equal(t, 3, v)

	v = none.UnwrapOrElse(func() int { return 2 })
	Equal(t, 2, v)

	v = none.UnwrapOrDefault()
	Equal(t, 0, v)

	// now test with a pointer type.
	type testStruct struct {
		S string
	}

	sNone := None[*testStruct]()
	PanicMatches(t, func() { sNone.Unwrap() }, "Option.Unwrap: option is None")

	v2 := sNone.UnwrapOr(&testStruct{S: "blah"})
	Equal(t, &testStruct{S: "blah"}, v2)

	v2 = sNone.UnwrapOrElse(func() *testStruct { return &testStruct{S: "blah 2"} })
	Equal(t, &testStruct{S: "blah 2"}, v2)

	v2 = sNone.UnwrapOrDefault()
	Equal(t, nil, v2)
}

func TestNilOption(t *testing.T) {
	value := Some[any](nil)
	Equal(t, false, value.IsNone())
	Equal(t, true, value.IsSome())
	Equal(t, nil, value.Unwrap())

	ret := returnTypedNoneOption()
	Equal(t, true, ret.IsNone())
	Equal(t, false, ret.IsSome())
	PanicMatches(t, func() {
		ret.Unwrap()
	}, "Option.Unwrap: option is None")

	ret = returnTypedSomeOption()
	Equal(t, false, ret.IsNone())
	Equal(t, true, ret.IsSome())
	Equal(t, testStruct{}, ret.Unwrap())

	retPtr := returnTypedNoneOptionPtr()
	Equal(t, true, retPtr.IsNone())
	Equal(t, false, retPtr.IsSome())

	retPtr = returnTypedSomeOptionPtr()
	Equal(t, false, retPtr.IsNone())
	Equal(t, true, retPtr.IsSome())
	Equal(t, new(testStruct), retPtr.Unwrap())
}

func TestOptionJSON(t *testing.T) {
	type s struct {
		Timestamp Option[time.Time] `json:"ts"`
	}

	now := time.Now().UTC().Truncate(time.Minute)
	tv := s{Timestamp: Some(now)}

	b, err := json.Marshal(tv)
	Equal(t, nil, err)
	Equal(t, `{"ts":"`+now.Format(time.RFC3339)+`"}`, string(b))

	tv = s{}
	b, err = json.Marshal(tv)
	Equal(t, nil, err)
	Equal(t, `{"ts":null}`, string(b))
}

func TestOptionJSONOmitempty(t *testing.T) {
	type s struct {
		Timestamp Option[time.Time] `json:"ts,omitempty"`
	}

	now := time.Now().UTC().Truncate(time.Minute)
	tv := s{Timestamp: Some(now)}

	b, err := json.Marshal(tv)
	Equal(t, nil, err)
	Equal(t, `{"ts":"`+now.Format(time.RFC3339)+`"}`, string(b))

	type s2 struct {
		Timestamp *Option[time.Time] `json:"ts,omitempty"`
	}
	tv2 := &s2{}
	b, err = json.Marshal(tv2)
	Equal(t, nil, err)
	Equal(t, `{}`, string(b))
}

func returnTypedNoneOption() Option[testStruct] {
	return None[testStruct]()
}

func returnTypedSomeOption() Option[testStruct] {
	return Some(testStruct{})
}

func returnTypedNoneOptionPtr() Option[*testStruct] {
	return None[*testStruct]()
}

func returnTypedSomeOptionPtr() Option[*testStruct] {
	return Some(new(testStruct))
}

func returnTypedSomeOptionNil() Option[any] {
	return Some[any](nil)
}

func returnTypedNoOption() *testStruct {
	return new(testStruct)
}

func returnNoOptionNil() (any, bool) {
	return nil, true
}
