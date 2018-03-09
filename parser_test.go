package colorjson

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	var tokens []Token
	require.NoError(t, parseValue(strings.NewReader(`{"a":[1,null],"b":true}`), func(T Token) {
		tokens = append(tokens, T)
	}))
	require.EqualValues(t, []Token{
		{'{', ObjectType},
		{'"', StringType},
		{'a', StringType},
		{'"', StringType},
		{':', ObjectType},
		{'[', ArrayType},
		{'1', NumberType},
		{',', ArrayType},
		{'n', NullType},
		{'u', NullType},
		{'l', NullType},
		{'l', NullType},
		{']', ArrayType},
		{',', ObjectType},
		{'"', StringType},
		{'b', StringType},
		{'"', StringType},
		{':', ObjectType},
		{'t', BoolType},
		{'r', BoolType},
		{'u', BoolType},
		{'e', BoolType},
		{'}', ObjectType},
	}, tokens)
}
