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
