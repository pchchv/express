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
