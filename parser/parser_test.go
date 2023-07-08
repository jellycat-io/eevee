package parser

import (
	"fmt"
	"testing"

	"github.com/jellycat-io/eevee/ast"
	"github.com/jellycat-io/eevee/lexer"
	"github.com/jellycat-io/eevee/test"
)

func TestParseProgram(t *testing.T) {
	input := test.MakeInput(
		`42`,
		`"eevee"`,
		`3.14`,
	)

	l := lexer.New(input, 4)
	fmt.Println("tokens", l.Tokens)
	p := New(l.Tokens, true)
	ast := p.Parse()

	expectedAst := makeProgram(
		makeExpressionStatement(makeIntegerLiteral(42)),
		makeExpressionStatement(makeStringLiteral("eevee")),
		makeExpressionStatement(makeFloatLiteral(3.14)),
	)

	if ast.String() != expectedAst.String() {
		t.Fatalf("Expected: %q, got %q", expectedAst, ast)
	}
}

func TestParseBlock(t *testing.T) {
	input := test.MakeInput(
		`42`,
		`	"eevee"`,
		`		3.14`,
		`"flareon"`,
	)

	l := lexer.New(input, 4)
	p := New(l.Tokens, true)
	ast := p.Parse()

	expectedAst := makeProgram(
		makeExpressionStatement(makeIntegerLiteral(42)),
		makeBlockStatement(
			makeExpressionStatement(makeStringLiteral("eevee")),
			makeBlockStatement(
				makeExpressionStatement(makeFloatLiteral(3.14)),
			),
		),
		makeExpressionStatement(makeStringLiteral("flareon")),
	)

	if ast.String() != expectedAst.String() {
		t.Fatalf("Expected: %q, got %q", expectedAst, ast)
	}
}

func TestParseAssignmentExpression(t *testing.T) {
	input := test.MakeInput(
		`pokemon = "eevee"`,
		`level += 1`,
		`pokemon = eevee`,
		`pokemon = eevee = flareon`,
		`pokemon = eevee = "eevee"`,
		`level = 40 + 2`,
	)

	l := lexer.New(input, 4)
	p := New(l.Tokens, false)
	ast := p.Parse()

	expectedAst := makeProgram(
		makeExpressionStatement(makeAssignmentExpression(
			"=",
			makeIdentifier("pokemon"),
			makeStringLiteral("eevee"),
		)),
		makeExpressionStatement(makeAssignmentExpression(
			"+=",
			makeIdentifier("level"),
			makeIntegerLiteral(1),
		)),
		makeExpressionStatement(makeAssignmentExpression(
			"=",
			makeIdentifier("pokemon"),
			makeIdentifier("eevee"),
		)),
		makeExpressionStatement(makeAssignmentExpression(
			"=",
			makeIdentifier("pokemon"),
			makeAssignmentExpression(
				"=",
				makeIdentifier("eevee"),
				makeIdentifier("flareon"),
			),
		)),
		makeExpressionStatement(makeAssignmentExpression(
			"=",
			makeIdentifier("pokemon"),
			makeAssignmentExpression(
				"=",
				makeIdentifier("eevee"),
				makeStringLiteral("eevee"),
			),
		)),
		makeExpressionStatement(makeAssignmentExpression(
			"=",
			makeIdentifier("level"),
			makeBinaryExpression(
				"+",
				makeIntegerLiteral(40),
				makeIntegerLiteral(2),
			),
		)),
	)

	if ast.String() != expectedAst.String() {
		t.Fatalf("Expected: %q, got %q", expectedAst, ast)
	}
}

func TestParseBinaryExpression(t *testing.T) {
	input := test.MakeInput(
		`2 + 2`,
		`2 - 2`,
		`2 * 2`,
		`2 / 2`,
		`2 % 2`,
		`2 + 2 * 2`,
		`2 * 2 + 2`,
		`2 * (2 + 2)`,
	)

	l := lexer.New(input, 4)
	p := New(l.Tokens, false)
	ast := p.Parse()

	expectedAst := makeProgram(
		makeExpressionStatement(makeBinaryExpression(
			"+",
			makeIntegerLiteral(2),
			makeIntegerLiteral(2),
		)),
		makeExpressionStatement(makeBinaryExpression(
			"-",
			makeIntegerLiteral(2),
			makeIntegerLiteral(2),
		)),
		makeExpressionStatement(makeBinaryExpression(
			"*",
			makeIntegerLiteral(2),
			makeIntegerLiteral(2),
		)),
		makeExpressionStatement(makeBinaryExpression(
			"/",
			makeIntegerLiteral(2),
			makeIntegerLiteral(2),
		)),
		makeExpressionStatement(makeBinaryExpression(
			"%",
			makeIntegerLiteral(2),
			makeIntegerLiteral(2),
		)),
		makeExpressionStatement(makeBinaryExpression(
			"+",
			makeIntegerLiteral(2),
			makeBinaryExpression(
				"*",
				makeIntegerLiteral(2),
				makeIntegerLiteral(2),
			),
		)),
		makeExpressionStatement(makeBinaryExpression(
			"+",
			makeBinaryExpression(
				"*",
				makeIntegerLiteral(2),
				makeIntegerLiteral(2),
			),
			makeIntegerLiteral(2),
		)),
		makeExpressionStatement(makeBinaryExpression(
			"*",
			makeIntegerLiteral(2),
			makeBinaryExpression(
				"+",
				makeIntegerLiteral(2),
				makeIntegerLiteral(2),
			),
		)),
	)

	if ast.String() != expectedAst.String() {
		t.Fatalf("Expected: %q, got %q", expectedAst, ast)
	}
}

func TestParseLiteral(t *testing.T) {
	input := test.MakeInput(
		`42`,
		`"eevee"`,
		`3.14`,
	)

	l := lexer.New(input, 4)
	p := New(l.Tokens, true)
	ast := p.Parse()

	expectedAst := makeProgram(
		makeExpressionStatement(makeIntegerLiteral(42)),
		makeExpressionStatement(makeStringLiteral("eevee")),
		makeExpressionStatement(makeFloatLiteral(3.14)),
	)

	if ast.String() != expectedAst.String() {
		t.Fatalf("Expected: %q, got %q", expectedAst, ast)
	}
}

func makeProgram(stmts ...ast.Statement) *ast.Program {
	s := []ast.Statement{}
	s = append(s, stmts...)
	return ast.NewProgram(s)
}

func makeBlockStatement(stmts ...ast.Statement) *ast.BlockStatement {
	s := []ast.Statement{}
	s = append(s, stmts...)
	return ast.NewBlockStatement(s)
}

func makeExpressionStatement(e ast.Expression) *ast.ExpressionStatement {
	return ast.NewExpressionStatement(e)
}

func makeAssignmentExpression(op string, l, r ast.Expression) *ast.AssignmentExpression {
	return ast.NewAssignmentExpression(op, l, r)
}

func makeBinaryExpression(op string, l, r ast.Expression) *ast.BinaryExpression {
	return ast.NewBinaryExpression(op, l, r)
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

func makeIdentifier(name string) *ast.Identifier {
	return ast.NewIdentifier(name)
}
