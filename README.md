# express [![CI](https://github.com/pchchv/express/workflows/CI/badge.svg)](https://github.com/pchchv/express/actions?query=workflow%3ACI+event%3Apush) [![Godoc Reference](https://pkg.go.dev/badge/github.com/pchchv/express)](https://pkg.go.dev/github.com/pchchv/express) [![Go Report Card](https://goreportcard.com/badge/github.com/pchchv/express)](https://goreportcard.com/report/github.com/pchchv/express)

Lexer, parser, cli, and library for working with JSON data expressions.

## Usage

```go
package main

import (
	"fmt"

	"github.com/pchchv/express"
)

func main() {
	expression := []byte(`.properties.employees > 20`)
	input := []byte(`{"name":"MyCompany", "properties":{"employees": 50}`)
	ex, err := express.Parse(expression)
	if err != nil {
		panic(err)
	}

	result, err := ex.Calculate(input)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", result)
}
```

## Expressions
Expressions support most mathematical and string expressions see below for details:

### Syntax & Rules

| Token          | Example                  | Syntax Rules                                                                                                                                                                              |
|----------------|--------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `Equals`       | `==`                     | supports both `==` and `=`.                                                                                                                                                               |
| `Add`          | `+`                      | N/A                                                                                                                                                                                       |
| `Subtract`     | `-`                      | N/A                                                                                                                                                                                       |
| `Multiply`     | `*`                      | N/A                                                                                                                                                                                       |
| `Divide`       | `/`                      | N/A                                                                                                                                                                                       |
| `Gt`           | `>`                      | N/A                                                                                                                                                                                       |
| `Gte`          | `>=`                     | N/A                                                                                                                                                                                       |
| `Lt`           | `<`                      | N/A                                                                                                                                                                                       |
| `Lte`          | `<=`                     | N/A                                                                                                                                                                                       |
| `OpenParen`    | `(`                      | N/A                                                                                                                                                                                       |
| `CloseParen`   | `)`                      | N/A                                                                                                                                                                                       |
| `OpenBracket`  | `[`                      | N/A                                                                                                                                                                                       |
| `CloseBracket` | `]`                      | N/A                                                                                                                                                                                       |
| `Comma`        | `,`                      | N/A                                                                                                                                                                                       |
| `QuotedString` | `"sample text"`          | Must start and end with an unescaped `"` character                                                                                                                                        |
| `Number`       | ` 123.45 `               | Must start and end with a space or '+' or '-' when hard coded value in expression and supports `0-9 +- e` characters for numbers and exponent notation.                                   |
| `BooleanTrue`  | `true`                   | Accepts `true` as a boolean only.                                                                                                                                                         |
| `BooleanFalse` | `false`                  | Accepts `false` as a boolean only.                                                                                                                                                        |
| `SelectorPath` | `.selector_path`         | Starts with a `.` and ends with whitespace blank space. This crate currently uses [gjson](https://github.com/tidwall/gjson.rs) and so the full gjson syntax for identifiers is supported. |
| `And`          | `&&`                     | N/A                                                                                                                                                                                       |
| `Not`          | `!`                      | Must be before Boolean identifier or expression or be followed by an operation                                                                                                            |
| `Or`           | <code>&vert;&vert;<code> | N/A                                                                                                                                                                                       |
| `Contains`     | `CONTAINS `              | Ends with whitespace blank space.                                                                                                                                                         |
| `ContainsAny`  | `CONTAINS_ANY `          | Ends with whitespace blank space.                                                                                                                                                         |
| `ContainsAll`  | `CONTAINS_ALL `          | Ends with whitespace blank space.                                                                                                                                                         |
| `In`           | `IN `                    | Ends with whitespace blank space.                                                                                                                                                         |
| `Between`      | ` BETWEEN `              | Starts & ends with whitespace blank space. example `1 BETWEEN 0 10`                                                                                                                       |
| `StartsWith`   | `STARTSWITH `            | Ends with whitespace blank space.                                                                                                                                                         |
| `EndsWith`     | `ENDSWITH `              | Ends with whitespace blank space.                                                                                                                                                         |
| `NULL`         | `NULL`                   | N/A                                                                                                                                                                                       |
| `Coerce`       | `COERCE`                 | Coerces one data type into another using in combination with 'Identifier'. Syntax is `COERCE <expression> _identifer_`.                                                                   |
| `Identifier`   | `_identifier_`           | Starts and end with an `_` used with 'COERCE' to cast data types, see table below with supported values. You can combine multiple coercions if separated by a COMMA.                      |
| `Colon`        | `:`                      | N/A                                                                                                                                                                                       |

### COERCE Types

| Type            | Description                                                                                                              |
|-----------------|--------------------------------------------------------------------------------------------------------------------------|
| `_datetime_`    | This attempts to convert the type into a DateTime.                                                                       |
| `_lowercase_`   | This converts the text into lowercase.                                                                                   |
| `_uppercase_`   | This converts the text into uppercase.                                                                                   |
| `_title_`       | This converts the text into title case, when the first letter is capitalized but the rest lower cased.                   |
| `_string_`      | This converts the value into a string and supports the Value's String, Number, Bool, DateTime with nanosecond precision. |
| `_number_`      | This converts the value into an f64 number and supports the Value's Null, String, Number, Bool and DateTime.             |
| `_substr_[n:n]` | This allows taking a substring of a string value. this returns Null if no match at specified indices exits.              |