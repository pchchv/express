package express

import "fmt"

// ErrUnterminatedString represents an unterminated string.
type ErrUnterminatedString struct {
	s string
}

func (e ErrUnterminatedString) Error() string {
	return fmt.Sprintf("Unterminated string `%s`", e.s)
}

// ErrInvalidSelectorPath represents an invalid selector string.
type ErrInvalidSelectorPath struct {
	s string
}

func (e ErrInvalidSelectorPath) Error() string {
	return fmt.Sprintf("Invalid selector path `%s`", e.s)
}

// ErrInvalidKeyword represents an invalid keyword keyword.
type ErrInvalidKeyword struct {
	s string
}

func (e ErrInvalidKeyword) Error() string {
	return fmt.Sprintf("Invalid keyword `%s`", e.s)
}

// ErrInvalidBool represents an invalid boolean.
type ErrInvalidBool struct {
	s string
}

func (e ErrInvalidBool) Error() string {
	return fmt.Sprintf("Invalid boolean `%s`", e.s)
}

// ErrInvalidNumber represents an invalid number.
type ErrInvalidNumber struct {
	s string
}

func (e ErrInvalidNumber) Error() string {
	return fmt.Sprintf("Invalid number `%s`", e.s)
}

// ErrUnsupportedCharacter represents an
// unsupported character is expression being lexed.
type ErrUnsupportedCharacter struct {
	b byte
}

func (e ErrUnsupportedCharacter) Error() string {
	return fmt.Sprintf("Unsupported Character `%s`", string(e.b))
}

// ErrInvalidIdentifier represents an invalid identifier.
type ErrInvalidIdentifier struct {
	s string
}

func (e ErrInvalidIdentifier) Error() string {
	return fmt.Sprintf("Invalid identifier `%s`", e.s)
}

// ErrInvalidCoerce represents an invalid Coerce error.
type ErrInvalidCoerce struct {
	Err error
}

func (e ErrInvalidCoerce) Error() string {
	return fmt.Sprintf("invalid COERCE: `%s`", e.Err.Error())
}
