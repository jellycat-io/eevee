package lexer

import (
	"testing"

	"github.com/jellycat-io/eevee/test"
	"github.com/jellycat-io/eevee/token"
)

func TestTokenizer(t *testing.T) {
	input := test.MakeInput(
		`#This is a comment`,
		`42`,
		`3.14`,
		`"eevee"`,
		`""`,
		`	"flareon"`,
		`"vaporeon"`,
	)

	expected := []token.Token{
		token.NewToken(token.INT, "42", 2, 1),
		token.NewToken(token.FLOAT, "3.14", 3, 1),
		token.NewToken(token.STRING, "\"eevee\"", 4, 1),
		token.NewToken(token.STRING, "\"\"", 5, 1),
		token.NewToken(token.INDENT, "", 6, 1),
		token.NewToken(token.STRING, "\"flareon\"", 6, 5),
		token.NewToken(token.DEDENT, "", 7, 1),
		token.NewToken(token.STRING, "\"vaporeon\"", 7, 1),
		token.NewToken(token.EOF, "", 8, 1),
	}

	l := NewLexer(input, 4)

	for i, tok := range expected {
		if tok != l.Tokens[i] {
			t.Fatalf("Tests[%d] - Wrong token. Expected = %q, got = %q", i, tok, l.Tokens[i])
		}
	}
}
