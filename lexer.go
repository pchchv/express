package express


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
