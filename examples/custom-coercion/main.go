package main

import (
	"fmt"
	"strings"

	"github.com/pchchv/express"
)

type Star struct {
	expression express.Expression
}

func (s *Star) Calculate(json []byte) (interface{}, error) {
	inner, err := s.expression.Calculate(json)
	if err != nil {
		return nil, err
	}

	switch t := inner.(type) {
	case string:
		return strings.Repeat("*", len(t)), nil
	default:
		return nil, fmt.Errorf("cannot star value %v", inner)
	}
}

func main() {
	// add custom coercion to the parser
	// coercions start and end with an _(underscore)
	guard := express.Coercions.Lock()
	guard.T["_star_"] = func(_ *express.Parser, constEligible bool, expression express.Expression) (stillConstEligible bool, e express.Expression, err error) {
		return constEligible, &Star{expression}, nil
	}
	guard.Unlock()

	expression := []byte(`COERCE "My Name" _star_`)
	input := []byte(`{}`)
	ex, err := express.Parse(expression)
	if err != nil {
		panic(err)
	}

	result, err := ex.Calculate(input)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", result)
}
