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
