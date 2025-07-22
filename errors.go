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
