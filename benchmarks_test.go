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
