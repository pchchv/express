package express


// Collect tokenizes the input and returns tokens or error lexing them.
func collect(src []byte) (tokens []Token, err error) {
	tokenizer := NewTokenizer(src)
	for {
		next := tokenizer.Next()
		if next.IsNone() {
			break
		}

		result := next.Unwrap()
		if result.IsErr() {
			return nil, result.Err()
		}

		tokens = append(tokens, result.Unwrap())
	}

	return
}
