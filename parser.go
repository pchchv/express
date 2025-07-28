package express

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pchchv/extender/resultext"
	"github.com/pchchv/goitertools"
)

var (
	_ Expression = (*eq)(nil)
	_ Expression = (*gt)(nil)
	_ Expression = (*or)(nil)
	_ Expression = (*lt)(nil)
	_ Expression = (*in)(nil)
	_ Expression = (*add)(nil)
	_ Expression = (*and)(nil)
	_ Expression = (*div)(nil)
	_ Expression = (*gte)(nil)
	_ Expression = (*num)(nil)
	_ Expression = (*lte)(nil)
	_ Expression = (*str)(nil)
	_ Expression = (*sub)(nil)
	_ Expression = (*not)(nil)
	_ Expression = (*null)(nil)
	_ Expression = (*array)(nil)
	_ Expression = (*multi)(nil)
	_ Expression = (*between)(nil)
	_ Expression = (*boolean)(nil)
	_ Expression = (*endsWith)(nil)
	_ Expression = (*contains)(nil)
	_ Expression = (*startsWith)(nil)
	_ Expression = (*containsAll)(nil)
	_ Expression = (*containsAny)(nil)
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

func (p *Parser) parseOperation(token Token, current Expression) (ex Expression, err error) {
	return
}

func (p *Parser) parseValue(token Token) (ex Expression, err error) {
	return
}

func (p *Parser) parseExpression() (current Expression, err error) {
	for {
		next := p.Tokenizer.Next()
		if next.IsNone() {
			return current, nil
		}

		result := next.Unwrap()
		if result.IsErr() {
			return nil, result.Err()
		}

		token := result.Unwrap()
		if current == nil {
			// look for nextToken value
			current, err = p.parseValue(token)
			if err != nil {
				return nil, err
			}
		} else {
			if token.Kind == CloseParen {
				return current, nil
			}

			// look for nextToken operation
			current, err = p.parseOperation(token, current)
			if err != nil {
				return nil, err
			}
		}
	}
}

func (p *Parser) nextOperatorToken(operationToken Token) (token Token, err error) {
	next := p.Tokenizer.Next()
	if next.IsNone() {
		start := int(operationToken.Start)
		err = fmt.Errorf("no value found after operation: %s", string(p.Exp[start:start+int(operationToken.Len)]))
		return
	}

	result := next.Unwrap()
	if result.IsErr() {
		err = result.Err()
		return
	}

	return result.Unwrap(), nil
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

type add struct {
	left  Expression
	right Expression
}

func (a add) Calculate(src []byte) (any, error) {
	left, err := a.left.Calculate(src)
	if err != nil {
		return nil, err
	}

	right, err := a.right.Calculate(src)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(left) != reflect.TypeOf(right) {
		if left != nil && right == nil {
			switch left.(type) {
			case string, float64:
				return left, nil
			}
		} else if right != nil && left == nil {
			switch right.(type) {
			case string, float64:
				return right, nil
			}
		}

		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s + %s", left, right)}
	}

	switch l := left.(type) {
	case string:
		return l + right.(string), nil
	case float64:
		return l + right.(float64), nil
	default:
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s + %s", left, right)}
	}
}

type endsWith struct {
	left  Expression
	right Expression
}

func (e endsWith) Calculate(src []byte) (any, error) {
	left, err := e.left.Calculate(src)
	if err != nil {
		return nil, err
	}

	right, err := e.right.Calculate(src)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(left) != reflect.TypeOf(right) {
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s ENDSWITH %s", left, right)}
	}

	switch l := left.(type) {
	case string:
		return strings.HasSuffix(l, right.(string)), nil
	default:
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s ENDSWITH %s !", left, right)}
	}
}

type sub struct {
	left  Expression
	right Expression
}

func (s sub) Calculate(src []byte) (any, error) {
	left, err := s.left.Calculate(src)
	if err != nil {
		return nil, err
	}

	right, err := s.right.Calculate(src)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(left) != reflect.TypeOf(right) {
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s - %s", left, right)}
	}

	switch l := left.(type) {
	case float64:
		return l - right.(float64), nil
	default:
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s - %s", left, right)}
	}
}

type multi struct {
	left  Expression
	right Expression
}

func (m multi) Calculate(src []byte) (any, error) {
	left, err := m.left.Calculate(src)
	if err != nil {
		return nil, err
	}

	right, err := m.right.Calculate(src)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(left) != reflect.TypeOf(right) {
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s * %s", left, right)}
	}

	switch l := left.(type) {
	case float64:
		return l * right.(float64), nil
	default:
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s * %s", left, right)}
	}
}

type div struct {
	left  Expression
	right Expression
}

func (d div) Calculate(src []byte) (any, error) {
	left, err := d.left.Calculate(src)
	if err != nil {
		return nil, err
	}

	right, err := d.right.Calculate(src)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(left) != reflect.TypeOf(right) {
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s / %s", left, right)}
	}

	switch l := left.(type) {
	case float64:
		return l / right.(float64), nil
	default:
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s / %s", left, right)}
	}
}

type eq struct {
	left  Expression
	right Expression
}

func (e eq) Calculate(src []byte) (any, error) {
	left, err := e.left.Calculate(src)
	if err != nil {
		return nil, err
	}

	right, err := e.right.Calculate(src)
	if err != nil {
		return nil, err
	}

	return reflect.DeepEqual(left, right), nil
}

type gt struct {
	left  Expression
	right Expression
}

func (g gt) Calculate(src []byte) (any, error) {
	left, err := g.left.Calculate(src)
	if err != nil {
		return nil, err
	}

	right, err := g.right.Calculate(src)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(left) != reflect.TypeOf(right) {
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s > %s", left, right)}
	}

	switch l := left.(type) {
	case string:
		return l > right.(string), nil
	case float64:
		return l > right.(float64), nil
	case time.Time:
		return l.After(right.(time.Time)), nil
	default:
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s > %s", left, right)}
	}
}

type gte struct {
	left  Expression
	right Expression
}

func (g gte) Calculate(src []byte) (any, error) {
	left, err := g.left.Calculate(src)
	if err != nil {
		return nil, err
	}

	right, err := g.right.Calculate(src)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(left) != reflect.TypeOf(right) {
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s >= %s", left, right)}
	}

	switch l := left.(type) {
	case string:
		return l >= right.(string), nil
	case float64:
		return l >= right.(float64), nil
	case time.Time:
		r := right.(time.Time)
		return l.After(r) || l.Equal(r), nil
	default:
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s >= %s", left, right)}
	}
}

type lt struct {
	left  Expression
	right Expression
}

func (l lt) Calculate(src []byte) (any, error) {
	left, err := l.left.Calculate(src)
	if err != nil {
		return nil, err
	}

	right, err := l.right.Calculate(src)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(left) != reflect.TypeOf(right) {
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s < %s", left, right)}
	}

	switch l := left.(type) {
	case string:
		return l < right.(string), nil
	case float64:
		return l < right.(float64), nil
	case time.Time:
		return l.Before(right.(time.Time)), nil
	default:
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s < %s", left, right)}
	}
}

type lte struct {
	left  Expression
	right Expression
}

func (l lte) Calculate(src []byte) (any, error) {
	left, err := l.left.Calculate(src)
	if err != nil {
		return nil, err
	}

	right, err := l.right.Calculate(src)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(left) != reflect.TypeOf(right) {
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s <= %s", left, right)}
	}

	switch l := left.(type) {
	case string:
		return l <= right.(string), nil
	case float64:
		return l <= right.(float64), nil
	case time.Time:
		r := right.(time.Time)
		return l.Before(r) || l.Equal(r), nil
	default:
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s <= %s", left, right)}
	}
}

type or struct {
	left  Expression
	right Expression
}

func (o or) Calculate(src []byte) (any, error) {
	left, err := o.left.Calculate(src)
	if err != nil {
		return nil, err
	}

	switch t := left.(type) {
	case bool:
		if t {
			return true, nil
		}
	}

	right, err := o.right.Calculate(src)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(left) != reflect.TypeOf(right) {
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s || %s", left, right)}
	}

	switch t := left.(type) {
	case bool:
		return t || right.(bool), nil
	default:
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s || %s !", left, right)}
	}
}

type and struct {
	left  Expression
	right Expression
}

func (a and) Calculate(src []byte) (any, error) {
	left, err := a.left.Calculate(src)
	if err != nil {
		return nil, err
	}

	switch t := left.(type) {
	case bool:
		if !t {
			return false, nil
		}
	default:
		return false, nil
	}

	right, err := a.right.Calculate(src)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(left) != reflect.TypeOf(right) {
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s && %s", left, right)}
	}

	switch t := left.(type) {
	case bool:
		return t && right.(bool), nil
	default:
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s && %s !", left, right)}
	}
}

type startsWith struct {
	left  Expression
	right Expression
}

func (s startsWith) Calculate(src []byte) (any, error) {
	left, err := s.left.Calculate(src)
	if err != nil {
		return nil, err
	}

	right, err := s.right.Calculate(src)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(left) != reflect.TypeOf(right) {
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s STARTSWITH %s", left, right)}
	}

	switch l := left.(type) {
	case string:
		return strings.HasPrefix(l, right.(string)), nil
	default:
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s STARTSWITH %s !", left, right)}
	}
}

type in struct {
	left  Expression
	right Expression
}

func (i in) Calculate(src []byte) (any, error) {
	left, err := i.left.Calculate(src)
	if err != nil {
		return nil, err
	}

	right, err := i.right.Calculate(src)
	if err != nil {
		return nil, err
	}

	arr, ok := right.([]any)
	if !ok {
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s IN %s !", left, right)}
	}

	for _, v := range arr {
		if left == v {
			return true, nil
		}
	}

	return false, nil
}

type contains struct {
	left  Expression
	right Expression
}

func (c contains) Calculate(src []byte) (any, error) {
	left, err := c.left.Calculate(src)
	if err != nil {
		return nil, err
	}

	right, err := c.right.Calculate(src)
	if err != nil {
		return nil, err
	}

	if leftTypeOf := reflect.TypeOf(left); leftTypeOf != reflect.TypeOf(right) && leftTypeOf.Kind() != reflect.Slice {
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s CONTAINS %s", left, right)}
	}

	switch l := left.(type) {
	case string:
		return strings.Contains(l, right.(string)), nil
	case []any:
		for _, v := range l {
			if reflect.DeepEqual(v, right) {
				return true, nil
			}
		}
		return false, nil
	default:
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s CONTAINS %s !", left, right)}
	}
}

type containsAny struct {
	left  Expression
	right Expression
}

func (c containsAny) Calculate(src []byte) (any, error) {
	left, err := c.left.Calculate(src)
	if err != nil {
		return nil, err
	}

	right, err := c.right.Calculate(src)
	if err != nil {
		return nil, err
	}

	switch l := left.(type) {
	case string:
		switch r := right.(type) {
		case string:
			for _, c := range r {
				for _, c2 := range l {
					if c == c2 {
						return true, nil
					}
				}
			}
		case []any:
			for _, v := range r {
				s, ok := v.(string)
				if !ok {
					continue
				}
				if strings.Contains(l, s) {
					return true, nil
				}
			}
			return false, nil
		default:
			return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s CONTAINS_ANY %s", left, right)}
		}
	case []any:
		switch r := right.(type) {
		case []any:
			for _, rv := range r {
				for _, lv := range l {
					if reflect.DeepEqual(rv, lv) {
						return true, nil
					}
				}
			}
		case string:
			for _, c := range r {
				for _, v := range l {
					if reflect.DeepEqual(string(c), v) {
						return true, nil
					}
				}
			}
		default:
			return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s CONTAINS_ANY %s", left, right)}
		}
	default:
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s CONTAINS_ANY %s !", left, right)}
	}
	return false, nil
}

type containsAll struct {
	left  Expression
	right Expression
}

func (c containsAll) Calculate(src []byte) (any, error) {
	left, err := c.left.Calculate(src)
	if err != nil {
		return nil, err
	}

	right, err := c.right.Calculate(src)
	if err != nil {
		return nil, err
	}

	switch l := left.(type) {
	case string:
		switch r := right.(type) {
		case string:
		OUTER1:
			for _, c := range r {
				for _, c2 := range l {
					if c == c2 {
						continue OUTER1
					}
				}
				return false, nil
			}
		case []any:
			for _, v := range r {
				s, ok := v.(string)
				if !ok || !strings.Contains(l, s) {
					return false, nil
				}
			}
			return true, nil
		default:
			return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s CONTAINS_ALL %s", left, right)}
		}
	case []any:
		switch r := right.(type) {
		case []any:
		OUTER3:
			for _, rv := range r {
				for _, lv := range l {
					if reflect.DeepEqual(rv, lv) {
						continue OUTER3
					}
				}
				return false, nil
			}
		case string:
		OUTER4:
			for _, c := range r {
				for _, v := range l {
					if reflect.DeepEqual(string(c), v) {
						continue OUTER4
					}
				}
				return false, nil
			}
		default:
			return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s CONTAINS_ALL %s", left, right)}
		}
	default:
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s CONTAINS_ALL %s !", left, right)}
	}
	return true, nil
}

type not struct {
	value Expression
}

func (n not) Calculate(src []byte) (any, error) {
	value, err := n.value.Calculate(src)
	if err != nil {
		return nil, err
	}

	switch t := value.(type) {
	case bool:
		return !t, nil
	default:
		return nil, ErrUnsupportedTypeComparison{s: fmt.Sprintf("%s for !", value)}
	}
}

type array struct {
	vec []Expression
}

func (a array) Calculate(src []byte) (any, error) {
	arr := make([]any, 0, len(a.vec))
	for _, v := range a.vec {
		res, err := v.Calculate(src)
		if err != nil {
			return nil, err
		}

		arr = append(arr, res)
	}

	return arr, nil
}

type num struct {
	n float64
}

func (n num) Calculate(_ []byte) (any, error) {
	return n.n, nil
}

type str struct {
	s string
}

func (s str) Calculate(_ []byte) (any, error) {
	return s.s, nil
}

type boolean struct {
	b bool
}

func (b boolean) Calculate(_ []byte) (any, error) {
	return b.b, nil
}

type null struct{}

func (bn null) Calculate(_ []byte) (any, error) {
	return nil, nil
}
