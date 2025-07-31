package express

import "testing"

func benchParsing(b *testing.B, expression string) {
	b.SetBytes(int64(len(expression)))
	for i := 0; i < b.N; i++ {
		if _, err := Parse([]byte(expression)); err != nil {
			b.Fatal(err)
		}
	}
}

func benchLexing(b *testing.B, expression string) {
	exp := []byte(expression)
	b.SetBytes(int64(len(expression)))
	for i := 0; i < b.N; i++ {
		if _, err := collect(exp); err != nil {
			b.Fatal(err)
		}
	}
}

func benchExecution(b *testing.B, expression, input string) {
	ex, err := Parse([]byte(expression))
	if err != nil {
		b.Fatal(err)
	}

	in := []byte(input)
	b.SetBytes(int64(len(in)))
	for i := 0; i < b.N; i++ {
		if _, err := ex.Calculate(in); err != nil {
			b.Fatal(err)
		}
	}
}
