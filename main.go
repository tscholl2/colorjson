package colorjson

import (
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
)

// TokenType represents the different types of JSON types.
type TokenType int

const (
	// ArrayType represents arrays.
	ArrayType = TokenType(iota)
	// ObjectType represents objects.
	ObjectType
	// NumberType represents numbers.
	NumberType
	// BoolType represents bools.
	BoolType
	// StringType represents strings.
	StringType
	// NullType represents nulls.
	NullType
)

// Token contains the current rune and what its type is.
type Token struct {
	Rune rune
	Type TokenType
}

// Print takes a string and prints it to the writer
// surrounding different json types by colors.
func Print(s *strings.Reader, w io.Writer) error {
	colors := []string{
		"\u001B[31m", // red
		"\u001B[34m", // blue
		"\u001B[32m", // green
		"\u001B[36m", // cyan
		"\u001B[33m", // yellow
		"\u001B[35m", // purple
	}
	return Parse(s, func(T Token) {
		fmt.Fprintf(w, colors[T.Type]+string(T.Rune)+"\u001B[0m")
	})
}

// Parse calls the given function on each token found.
func Parse(s *strings.Reader, fn func(T Token)) error {
	return parseValue(s, fn)
}

type action func(T Token)

func parseValue(reader *strings.Reader, fn action) error {
	switch r := peek(reader); {
	case r == '"':
		return parseString(reader, fn)
	case ('0' <= r && r <= '9') || r == '-':
		return parseNumber(reader, fn)
	case r == '[':
		return parseArray(reader, fn)
	case r == 't' || r == 'f':
		return parseBool(reader, fn)
	case r == 'n':
		return parseNull(reader, fn)
	case r == '{':
		return parseObject(reader, fn)
	default:
		return errors.New("unknown type")
	}
}

func parseObject(reader *strings.Reader, fn action) error {
	if err := parseLiteral(reader, ObjectType, fn, "{"); err != nil {
		return err
	}
	for r := peek(reader); r != '}'; r = peek(reader) {
		if err := parseWhitespace(reader, ObjectType, fn); err != nil {
			return err
		}
		if err := parseString(reader, fn); err != nil {
			return err
		}
		if err := parseWhitespace(reader, ObjectType, fn); err != nil {
			return err
		}
		if err := parseLiteral(reader, ObjectType, fn, ":"); err != nil {
			return err
		}
		if err := parseWhitespace(reader, ObjectType, fn); err != nil {
			return err
		}
		if err := parseValue(reader, fn); err != nil {
			return err
		}
		if peek(reader) == ',' {
			parseLiteral(reader, ObjectType, fn, ",")
		}
		if err := parseWhitespace(reader, ObjectType, fn); err != nil {
			return err
		}
	}
	return parseLiteral(reader, ObjectType, fn, "}")
}

func parseArray(reader *strings.Reader, fn action) error {
	if err := parseLiteral(reader, ArrayType, fn, "["); err != nil {
		return err
	}
	for r := peek(reader); r != ']'; r = peek(reader) {
		if err := parseWhitespace(reader, ArrayType, fn); err != nil {
			return err
		}
		if err := parseValue(reader, fn); err != nil {
			return err
		}
		if err := parseWhitespace(reader, ArrayType, fn); err != nil {
			return err
		}
		if peek(reader) == ',' {
			parseLiteral(reader, ArrayType, fn, ",")
		}
		if err := parseWhitespace(reader, ArrayType, fn); err != nil {
			return err
		}
	}
	return parseLiteral(reader, ArrayType, fn, "]")
}

func parseString(reader *strings.Reader, fn action) error {
	if err := parseLiteral(reader, StringType, fn, `"`); err != nil {
		return err
	}
	var r rune
	var err error
	for n := peek(reader); n != '"' || r == '\\'; n = peek(reader) {
		if r, _, err = reader.ReadRune(); err != nil {
			return errors.Wrapf(err, "str:")
		}
		fn(Token{r, StringType})
	}
	return parseLiteral(reader, StringType, fn, `"`)
}

func parseBool(reader *strings.Reader, fn action) error {
	if r := peek(reader); r != 't' && r != 'f' {
		return unexpected(string(r), "t, f")
	}
	if peek(reader) == 't' {
		return parseLiteral(reader, BoolType, fn, "true")
	}
	return parseLiteral(reader, BoolType, fn, "false")
}

func parseNull(reader *strings.Reader, fn action) error {
	return parseLiteral(reader, NullType, fn, "null")
}

func parseNumber(reader *strings.Reader, fn action) (err error) {
	var r rune
	for err == nil && strings.ContainsRune("0123456789.-eE", peek(reader)) {
		r, _, err = reader.ReadRune()
		fn(Token{r, NumberType})
	}
	return
}

func parseLiteral(reader *strings.Reader, t TokenType, fn action, literal string) error {
	b := make([]byte, len(literal))
	if _, err := reader.Read(b); err != nil {
		return err
	}
	if string(b) == literal {
		for _, r := range literal {
			fn(Token{r, t})
		}
		return nil
	}
	return unexpected(string(b), literal)
}

func parseWhitespace(reader *strings.Reader, t TokenType, fn action) error {
	for r := peek(reader); r == ' ' || r == '\n' || r == '\t'; r = peek(reader) {
		if _, _, err := reader.ReadRune(); err != nil {
			return err
		}
		fn(Token{r, t})
	}
	return nil
}

func unexpected(got string, expected string) error {
	return errors.Wrapf(fmt.Errorf("unexpected '%s', expected '%s", got, expected), "syntax")
}

func peek(reader *strings.Reader) rune {
	r, _, _ := reader.ReadRune()
	reader.UnreadRune()
	return r
}
