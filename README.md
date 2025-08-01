# express [![CI](https://github.com/pchchv/express/workflows/CI/badge.svg)](https://github.com/pchchv/express/actions?query=workflow%3ACI+event%3Apush) [![Godoc Reference](https://pkg.go.dev/badge/github.com/pchchv/express)](https://pkg.go.dev/github.com/pchchv/express) [![Go Report Card](https://goreportcard.com/badge/github.com/pchchv/express)](https://goreportcard.com/report/github.com/pchchv/express)

Lexer, parser, cli, and library for working with JSON data expressions.

## Usage

```go
package main

import (
	"fmt"

	"github.com/pchchv/express"
)

func main() {
	expression := []byte(`.properties.employees > 20`)
	input := []byte(`{"name":"MyCompany", "properties":{"employees": 50}`)
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
```
