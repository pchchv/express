package resultext

import "errors"

type Struct struct{}

func returnOk() Result[Struct, error] {
	return Ok[Struct, error](Struct{})
}

func returnErr() Result[Struct, error] {
	return Err[Struct, error](errors.New("bad"))
}
