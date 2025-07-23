package express

import (
	"github.com/pchchv/extender/resultext"
	"github.com/pchchv/goitertools"
)

// Expression Represents a stateless parsed expression that can be applied to JSON data.
type Expression interface {
	// Calculate executes the parsed expression and apply it against the supplied data.
	//
	// Will return `Err` if the expression cannot be applied to the
	// supplied data due to invalid data type comparisons.
	Calculate(src []byte) (any, error)
}

// Parser parses and returns a supplied expression.
type Parser struct {
	Exp       []byte
	Tokenizer goitertools.PeekableIterator[resultext.Result[Token, error]]
}
