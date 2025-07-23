package express

import (
	"fmt"
	"reflect"
	"time"

	"github.com/pchchv/extender/resultext"
	"github.com/pchchv/goitertools"
)

var _ Expression = (*between)(nil)

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

type between struct {
	left  Expression
	right Expression
	value Expression
}

func (b between) Calculate(src []byte) (any, error) {
	left, err := b.left.Calculate(src)
	if err != nil {
		return nil, err
	}

	right, err := b.right.Calculate(src)
	if err != nil {
		return nil, err
	}

	value, err := b.value.Calculate(src)
	if err != nil {
		return nil, err
	}

	// fast path, if any are nil/null there's no way to actually do the BETWEEN comparison
	if left == nil || right == nil || value == nil {
		return false, nil
	}

	leftType := reflect.TypeOf(left)
	if !(leftType == reflect.TypeOf(right) && reflect.TypeOf(value) == leftType) {
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s < %s", left, right)}
	}

	switch v := value.(type) {
	case string:
		return v > left.(string) && v < right.(string), nil
	case float64:
		return v > left.(float64) && v < right.(float64), nil
	case time.Time:
		return v.After(left.(time.Time)) && v.Before(right.(time.Time)), nil
	default:
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s < %s", left, right)}
	}
}

type endsWith struct {
	left  Expression
	right Expression
}
