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
