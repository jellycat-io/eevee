package parser

import (
	"testing"

	"github.com/jellycat-io/eevee/ast"
	"github.com/jellycat-io/eevee/test"
)

func TestParseProgram(t *testing.T) {
	input := test.MakeInput(
		"INT=42",
		`STRING="eevee"`,
		`FLOAT=3.14`,
	)

	tokens := test.TokensFromString(input)
	p := NewParser(tokens)
	ast := p.Parse()

	expectedAst := makeProgram(
		makeExpressionStatement(makeIntegerLiteral(42)),
		makeExpressionStatement(makeStringLiteral("eevee")),
		makeExpressionStatement(makeFloatLiteral(3.14)),
	)

	if ast.String() != expectedAst.String() {
		t.Fatalf("Expected: %q, got %q", ast, expectedAst)
	}
}

func makeProgram(stmts ...ast.Statement) *ast.Program {
	s := []ast.Statement{}
	s = append(s, stmts...)
	return ast.NewProgram(s)
}

func makeExpressionStatement(e ast.Expression) *ast.ExpressionStatement {
	return ast.NewExpressionStatement(e)
}

func makeIntegerLiteral(n int64) *ast.IntegerLiteral {
	return ast.NewIntegerLiteral(n)
}

func makeFloatLiteral(n float64) *ast.FloatLiteral {
	return ast.NewFloatLiteral(n)
}

func makeStringLiteral(s string) *ast.StringLiteral {
	return ast.NewStringLiteral(s)
}
