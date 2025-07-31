package express

import "testing"

func BenchmarkParsingNumPlusNum(b *testing.B) {
	benchParsing(b, "1 + 1")
}

func BenchmarkParsingIdentNum(b *testing.B) {
	benchParsing(b, ".field1 + 1")
}

func BenchmarkParsingIdentIdent(b *testing.B) {
	benchParsing(b, ".field1 + .field2")
}

func BenchmarkParsingFNameLName(b *testing.B) {
	benchParsing(b, `.first_name + " " + .last_name`)
}

func BenchmarkParsingParenDiv(b *testing.B) {
	benchParsing(b, `(1 + 1) / 2`)
}

func BenchmarkParsingParenDivIdents(b *testing.B) {
	benchParsing(b, `(.field1 + .field2) / .field3`)
}

func BenchmarkParsingCompanyEmployees(b *testing.B) {
	benchParsing(b, `.properties.employees > 20`)
}

func BenchmarkParsingParenNot(b *testing.B) {
	benchParsing(b, `!(.f1 != .f2)`)
}

func BenchmarkParsingCoerceDateTimeSelector(b *testing.B) {
	benchParsing(b, `COERCE .dt1 _datetime_ == COERCE .dt2 _datetime_`)
}

func BenchmarkParsingCoerceDateTimeSelectorMixed(b *testing.B) {
	benchParsing(b, `COERCE .dt1 _datetime_ == COERCE "2022-01-02" _datetime_`)
}

func BenchmarkParsingCoerceDateTimeSelectorConstant(b *testing.B) {
	benchParsing(b, `COERCE "2022-01-02" _datetime_ == COERCE "2022-01-02" _datetime_`)
}

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
