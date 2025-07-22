package express

const (
	SelectorPath = iota
	QuotedString
	Number
	BooleanTrue
	BooleanFalse
	Null
	Equals
	Add
	Subtract
	Multiply
	Divide
	Gt
	Gte
	Lt
	Lte
	And
	Or
	Not
	Contains
	ContainsAny
	ContainsAll
	In
	Between
	StartsWith
	EndsWith
	OpenBracket
	CloseBracket
	Comma
	OpenParen
	CloseParen
	Coerce
	Identifier
	Colon
)

// TokenKind is the type of token lexed.
type TokenKind uint8

// Token represents a lexed token.
type Token struct {
	Start uint32
	Len   uint16
	Kind  TokenKind
}

// LexerResult represents a token lexed result.
type LexerResult struct {
	kind TokenKind
	len  uint16
}

func isUpper(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func isLower(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isWhitespace(b byte) bool {
	switch b {
	case '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0:
		return true
	default:
		return false
	}
}

func isAlphanumeric(c byte) bool {
	return isLower(c) || isUpper(c) || isDigit(c)
}

func isAlphabetical(c byte) bool {
	return isLower(c) || isUpper(c)
}

func skipWhitespace(data []byte) uint16 {
	return takeWhile(data, func(b byte) bool {
		return isWhitespace(b)
	})
}

// takeWhile сonsumes bytes while a predicate evaluates to true.
func takeWhile(data []byte, pred func(byte) bool) (end uint16) {
	for _, b := range data {
		if !pred(b) {
			break
		}
		end++
	}
	return
}

func tokenizeSelectorPath(data []byte) (result LexerResult, err error) {
	if end := takeWhile(data[1:], func(b byte) bool {
		return !isWhitespace(b) && b != ')' && b != ']'
	}); end > 0 {
		if len(data) > int(end) {
			end += 1
		}
		result = LexerResult{
			kind: SelectorPath,
			len:  end,
		}
	} else {
		err = ErrInvalidSelectorPath{s: string(data)}
	}

	return
}

func tokenizeString(data []byte, quote byte) (result LexerResult, err error) {
	var lastBackslash, endedWithTerminator bool
	if end := takeWhile(data[1:], func(b byte) bool {
		switch b {
		case '\\':
			lastBackslash = true
			return true
		case quote:
			if lastBackslash {
				lastBackslash = false
				return true
			}
			endedWithTerminator = true
			return false
		default:
			return true
		}
	}); end > 0 {
		if endedWithTerminator {
			result = LexerResult{
				kind: QuotedString,
				len:  end + 2,
			}
		} else {
			err = ErrUnterminatedString{s: string(data)}
		}
	} else {
		if !endedWithTerminator || len(data) < 2 {
			err = ErrUnterminatedString{s: string(data)}
		} else {
			result = LexerResult{
				kind: QuotedString,
				len:  2,
			}
		}
	}

	return
}

func tokenizeNull(data []byte) (result LexerResult, err error) {
	if end := takeWhile(data, func(b byte) bool {
		return isAlphabetical(b)
	}); end > 0 && string(data[:end]) == "NULL" {
		result = LexerResult{
			kind: Null,
			len:  end,
		}
	} else {
		err = ErrInvalidKeyword{s: string(data)}
	}
	return
}

func tokenizeBool(data []byte) (result LexerResult, err error) {
	if end := takeWhile(data, func(b byte) bool {
		return isAlphabetical(b)
	}); end > 0 {
		switch string(data[:end]) {
		case "true":
			result = LexerResult{
				kind: BooleanTrue,
				len:  end,
			}
		case "false":
			result = LexerResult{
				kind: BooleanFalse,
				len:  end,
			}
		default:
			err = ErrInvalidBool{s: string(data)}
		}
	} else {
		err = ErrInvalidBool{s: string(data)}
	}
	return
}
