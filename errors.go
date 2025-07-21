package express

import "fmt"

// ErrUnterminatedString represents an unterminated string.
type ErrUnterminatedString struct {
	s string
}

func (e ErrUnterminatedString) Error() string {
	return fmt.Sprintf("Unterminated string `%s`", e.s)
}
