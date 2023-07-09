package parser

import (
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
	p := New(l.Tokens, true)
	ast := p.Parse()

	checkParserErrors(t, p)

	expectedAst := makeProgram(
		makeExpressionStatement(makeIntegerLiteral(42)),
		makeExpressionStatement(makeStringLiteral("eevee")),
		makeExpressionStatement(makeFloatLiteral(3.14)),
	)

	if ast.String() != expectedAst.String() {
		t.Fatalf("Expected: %q, got %q", expectedAst, ast)
	}
}

func TestParseBlockStatement(t *testing.T) {
	input := test.MakeInput(
		`42`,
		`	"eevee"`,
		`		3.14`,
		`"flareon"`,
	)

	l := lexer.New(input, 4)
	p := New(l.Tokens, true)
	ast := p.Parse()

	checkParserErrors(t, p)

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

func TestParseFunctionDeclaration(t *testing.T) {
	input := test.MakeInput(
		`fn square(x)`,
		`	return x * x`,
		`fn add(x, y) return x + y`,
		`fn nothing() return`,
	)

	l := lexer.New(input, 4)
	p := New(l.Tokens, true)
	ast := p.Parse()

	checkParserErrors(t, p)

	expectedAst := makeProgram(
		makeFunctionDeclaration(
			*makeIdentifier("square"),
			makeFunctionParameters(
				*makeIdentifier("x"),
			),
			makeBlockStatement(
				makeReturnStatement(
					makeBinaryExpression(
						"*",
						makeIdentifier("x"),
						makeIdentifier("x"),
					),
				),
			),
		),
		makeFunctionDeclaration(
			*makeIdentifier("add"),
			makeFunctionParameters(
				*makeIdentifier("x"),
				*makeIdentifier("y"),
			),
			makeReturnStatement(
				makeBinaryExpression(
					"+",
					makeIdentifier("x"),
					makeIdentifier("y"),
				),
			),
		),
		makeFunctionDeclaration(
			*makeIdentifier("nothing"),
			makeFunctionParameters(),
			makeReturnStatement(
				makeNullLiteral(),
			),
		),
	)

	if ast.String() != expectedAst.String() {
		t.Fatalf("Expected: %q, got %q", expectedAst, ast)
	}
}

func TestParseVariableStatement(t *testing.T) {
	input := test.MakeInput(
		`let pokemon = "eevee"`,
		`	let pokemon = eevee`,
		`let x, y`,
		`let x, y = 42`,
		`let x = 40 + 2`,
		`let x = y = 42`,
	)

	l := lexer.New(input, 4)
	p := New(l.Tokens, true)
	ast := p.Parse()

	checkParserErrors(t, p)

	expectedAst := makeProgram(
		makeVariableStatement(makeVariableDeclaration(
			makeIdentifier("pokemon"),
			makeStringLiteral("eevee"),
		)),
		makeBlockStatement(makeVariableStatement(makeVariableDeclaration(
			makeIdentifier("pokemon"),
			makeIdentifier("eevee"),
		))),
		makeVariableStatement(
			makeVariableDeclaration(
				makeIdentifier("x"),
				nil,
			),
			makeVariableDeclaration(
				makeIdentifier("y"),
				nil,
			),
		),
		makeVariableStatement(
			makeVariableDeclaration(
				makeIdentifier("x"),
				nil,
			),
			makeVariableDeclaration(
				makeIdentifier("y"),
				makeIntegerLiteral(42),
			),
		),
		makeVariableStatement(makeVariableDeclaration(
			makeIdentifier("x"),
			makeBinaryExpression(
				"+",
				makeIntegerLiteral(40),
				makeIntegerLiteral(2),
			),
		)),
		makeVariableStatement(makeVariableDeclaration(
			makeIdentifier("x"),
			makeAssignmentExpression(
				"=",
				makeIdentifier("y"),
				makeIntegerLiteral(42),
			),
		)),
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

	checkParserErrors(t, p)

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

func TestParseWhileStatement(t *testing.T) {
	input := test.MakeInput(
		`while x < 10 do`,
		`	x += 1`,
		`while true do x += 2`,
		`do x += 1 while x < 10`,
	)

	l := lexer.New(input, 4)
	p := New(l.Tokens, false)
	ast := p.Parse()

	checkParserErrors(t, p)

	expectedAst := makeProgram(
		makeWhileStatement(
			makeBinaryExpression(
				"<",
				makeIdentifier("x"),
				makeIntegerLiteral(10),
			),
			makeBlockStatement(
				makeExpressionStatement(
					makeAssignmentExpression(
						"+=",
						makeIdentifier("x"),
						makeIntegerLiteral(1),
					),
				),
			),
		),
		makeWhileStatement(
			makeBoolLiteral(true),
			makeExpressionStatement(
				makeAssignmentExpression(
					"+=",
					makeIdentifier("x"),
					makeIntegerLiteral(2),
				),
			),
		),
		makeDoWhileStatement(
			makeBinaryExpression(
				"<",
				makeIdentifier("x"),
				makeIntegerLiteral(10),
			),
			makeExpressionStatement(
				makeAssignmentExpression(
					"+=",
					makeIdentifier("x"),
					makeIntegerLiteral(1),
				),
			),
		),
	)

	if ast.String() != expectedAst.String() {
		t.Fatalf("Expected: %q, got %q", expectedAst, ast)
	}
}

func TestParseForStatement(t *testing.T) {
	input := test.MakeInput(
		`for let x = 1; x < 10; x += 1 do`,
		`	y += 1`,
		`for x = 1; x < 10; x += 1 do y += 1`,
		`for ;; do y += 1`,
	)

	l := lexer.New(input, 4)
	p := New(l.Tokens, false)
	ast := p.Parse()

	checkParserErrors(t, p)

	expectedAst := makeProgram(
		makeForStatement(
			makeVariableStatement(
				makeVariableDeclaration(
					makeIdentifier("x"),
					makeIntegerLiteral(1),
				),
			),
			makeBinaryExpression(
				"<",
				makeIdentifier("x"),
				makeIntegerLiteral(10),
			),
			makeAssignmentExpression(
				"+=",
				makeIdentifier("x"),
				makeIntegerLiteral(1),
			),
			makeBlockStatement(
				makeExpressionStatement(
					makeAssignmentExpression(
						"+=",
						makeIdentifier("y"),
						makeIntegerLiteral(1),
					),
				),
			),
		),
		makeForStatement(
			makeAssignmentExpression(
				"=",
				makeIdentifier("x"),
				makeIntegerLiteral(1),
			),
			makeBinaryExpression(
				"<",
				makeIdentifier("x"),
				makeIntegerLiteral(10),
			),
			makeAssignmentExpression(
				"+=",
				makeIdentifier("x"),
				makeIntegerLiteral(1),
			),
			makeExpressionStatement(
				makeAssignmentExpression(
					"+=",
					makeIdentifier("y"),
					makeIntegerLiteral(1),
				),
			),
		),
		makeForStatement(
			nil,
			nil,
			nil,
			makeExpressionStatement(
				makeAssignmentExpression(
					"+=",
					makeIdentifier("y"),
					makeIntegerLiteral(1),
				),
			),
		),
	)

	if ast.String() != expectedAst.String() {
		t.Fatalf("Expected: %q, got %q", expectedAst, ast)
	}
}

func TestParseIfStatement(t *testing.T) {
	input := test.MakeInput(
		`if level >= 15 == true then`,
		`	pokemon = "ivysaur"`,
		`else`,
		`	pokemon = "bulbasaur"`,
		`if (eevee not null) then`,
		`	if evo_cond is solar_stone then`,
		`		eevee = "leafeon"`,
		`	if evo_cond == friendship_plus_exchange then eevee = "sylveon"`,
		`	if evo_cond == friendship_at_night then eevee = "umbreon" else eevee = "espeon"`,
		`else eevee = "missingno"`,
	)

	l := lexer.New(input, 4)
	p := New(l.Tokens, false)
	ast := p.Parse()

	checkParserErrors(t, p)

	expectedAst := makeProgram(
		makeIfStatement(
			makeBinaryExpression(
				"==",
				makeBinaryExpression(
					">=",
					makeIdentifier("level"),
					makeIntegerLiteral(15),
				),
				makeBoolLiteral(true),
			),
			makeBlockStatement(
				makeExpressionStatement(
					makeAssignmentExpression(
						"=",
						makeIdentifier("pokemon"),
						makeStringLiteral("ivysaur"),
					),
				),
			),
			makeBlockStatement(
				makeExpressionStatement(
					makeAssignmentExpression(
						"=",
						makeIdentifier("pokemon"),
						makeStringLiteral("bulbasaur"),
					),
				),
			),
		),
		makeIfStatement(
			makeBinaryExpression(
				"!=",
				makeIdentifier("eevee"),
				makeNullLiteral(),
			),
			makeBlockStatement(
				makeIfStatement(
					makeBinaryExpression(
						"==",
						makeIdentifier("evo_cond"),
						makeIdentifier("solar_stone"),
					),
					makeBlockStatement(
						makeExpressionStatement(
							makeAssignmentExpression(
								"=",
								makeIdentifier("eevee"),
								makeStringLiteral("leafeon"),
							),
						),
					),
					nil,
				),
				makeIfStatement(
					makeBinaryExpression(
						"==",
						makeIdentifier("evo_cond"),
						makeIdentifier("friendship_plus_exchange"),
					),
					makeExpressionStatement(
						makeAssignmentExpression(
							"=",
							makeIdentifier("eevee"),
							makeStringLiteral("sylveon"),
						),
					),
					nil,
				),
				makeIfStatement(
					makeBinaryExpression(
						"==",
						makeIdentifier("evo_cond"),
						makeIdentifier("friendship_at_night"),
					),
					makeExpressionStatement(
						makeAssignmentExpression(
							"=",
							makeIdentifier("eevee"),
							makeStringLiteral("umbreon"),
						),
					),
					makeExpressionStatement(
						makeAssignmentExpression(
							"=",
							makeIdentifier("eevee"),
							makeStringLiteral("espeon"),
						),
					),
				),
			),
			makeExpressionStatement(
				makeAssignmentExpression(
					"=",
					makeIdentifier("eevee"),
					makeStringLiteral("missingno"),
				),
			),
		),
	)

	if ast.String() != expectedAst.String() {
		t.Fatalf("Expected: %q, got %q", expectedAst, ast)
	}
}

func TestParseLogicalExpression(t *testing.T) {
	input := test.MakeInput(
		`5 == 5 and 5 < 10`,
		`5 == 5 or 5 < 10`,
		`(5 == 5 && 5 < 10) and 5 > 1`,
		`(5 == 5 and 5 < 10) || 5 > 1`,
		`(5 == 5 or 5 < 10) || 5 > 1`,
		`(5 == 5 || 5 < 10) && 5 > 1`,
	)

	l := lexer.New(input, 4)
	p := New(l.Tokens, false)
	ast := p.Parse()

	checkParserErrors(t, p)

	expectedAst := makeProgram(
		makeExpressionStatement(makeLogicalExpression(
			"&&",
			makeBinaryExpression(
				"==",
				makeIntegerLiteral(5),
				makeIntegerLiteral(5),
			),
			makeBinaryExpression(
				"<",
				makeIntegerLiteral(5),
				makeIntegerLiteral(10),
			),
		)),
		makeExpressionStatement(makeLogicalExpression(
			"||",
			makeBinaryExpression(
				"==",
				makeIntegerLiteral(5),
				makeIntegerLiteral(5),
			),
			makeBinaryExpression(
				"<",
				makeIntegerLiteral(5),
				makeIntegerLiteral(10),
			),
		)),
		makeExpressionStatement(makeLogicalExpression(
			"&&",
			makeLogicalExpression(
				"&&",
				makeBinaryExpression(
					"==",
					makeIntegerLiteral(5),
					makeIntegerLiteral(5),
				),
				makeBinaryExpression(
					"<",
					makeIntegerLiteral(5),
					makeIntegerLiteral(10),
				),
			),
			makeBinaryExpression(
				">",
				makeIntegerLiteral(5),
				makeIntegerLiteral(1),
			),
		)),
		makeExpressionStatement(makeLogicalExpression(
			"||",
			makeLogicalExpression(
				"&&",
				makeBinaryExpression(
					"==",
					makeIntegerLiteral(5),
					makeIntegerLiteral(5),
				),
				makeBinaryExpression(
					"<",
					makeIntegerLiteral(5),
					makeIntegerLiteral(10),
				),
			),
			makeBinaryExpression(
				">",
				makeIntegerLiteral(5),
				makeIntegerLiteral(1),
			),
		)),
		makeExpressionStatement(makeLogicalExpression(
			"||",
			makeLogicalExpression(
				"||",
				makeBinaryExpression(
					"==",
					makeIntegerLiteral(5),
					makeIntegerLiteral(5),
				),
				makeBinaryExpression(
					"<",
					makeIntegerLiteral(5),
					makeIntegerLiteral(10),
				),
			),
			makeBinaryExpression(
				">",
				makeIntegerLiteral(5),
				makeIntegerLiteral(1),
			),
		)),
		makeExpressionStatement(makeLogicalExpression(
			"&&",
			makeLogicalExpression(
				"||",
				makeBinaryExpression(
					"==",
					makeIntegerLiteral(5),
					makeIntegerLiteral(5),
				),
				makeBinaryExpression(
					"<",
					makeIntegerLiteral(5),
					makeIntegerLiteral(10),
				),
			),
			makeBinaryExpression(
				">",
				makeIntegerLiteral(5),
				makeIntegerLiteral(1),
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
		`2 > 2`,
		`2 >= 2`,
		`2 < 2`,
		`2 <= 2`,
		`2 < 2 + 2`,
		`x = 2 > 2`,
		`2 == 2`,
		`2 is 2`,
		`4 != 2`,
		`4 not 2`,
		`2 not 2 < 2`,
		`2 == 2 < 2 + 2`,
		`-2 + 2`,
	)

	l := lexer.New(input, 4)
	p := New(l.Tokens, false)
	ast := p.Parse()

	checkParserErrors(t, p)

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
		makeExpressionStatement(makeBinaryExpression(
			">",
			makeIntegerLiteral(2),
			makeIntegerLiteral(2),
		)),
		makeExpressionStatement(makeBinaryExpression(
			">=",
			makeIntegerLiteral(2),
			makeIntegerLiteral(2),
		)),
		makeExpressionStatement(makeBinaryExpression(
			"<",
			makeIntegerLiteral(2),
			makeIntegerLiteral(2),
		)),
		makeExpressionStatement(makeBinaryExpression(
			"<=",
			makeIntegerLiteral(2),
			makeIntegerLiteral(2),
		)),
		makeExpressionStatement(makeBinaryExpression(
			"<",
			makeIntegerLiteral(2),
			makeBinaryExpression(
				"+",
				makeIntegerLiteral(2),
				makeIntegerLiteral(2),
			),
		)),
		makeExpressionStatement(makeAssignmentExpression(
			"=",
			makeIdentifier("x"),
			makeBinaryExpression(
				">",
				makeIntegerLiteral(2),
				makeIntegerLiteral(2),
			),
		)),
		makeExpressionStatement(makeBinaryExpression(
			"==",
			makeIntegerLiteral(2),
			makeIntegerLiteral(2),
		)),
		makeExpressionStatement(makeBinaryExpression(
			"==",
			makeIntegerLiteral(2),
			makeIntegerLiteral(2),
		)),
		makeExpressionStatement(makeBinaryExpression(
			"!=",
			makeIntegerLiteral(4),
			makeIntegerLiteral(2),
		)),
		makeExpressionStatement(makeBinaryExpression(
			"!=",
			makeIntegerLiteral(4),
			makeIntegerLiteral(2),
		)),
		makeExpressionStatement(makeBinaryExpression(
			"!=",
			makeIntegerLiteral(2),
			makeBinaryExpression(
				"<",
				makeIntegerLiteral(2),
				makeIntegerLiteral(2),
			),
		)),
		makeExpressionStatement(makeBinaryExpression(
			"==",
			makeIntegerLiteral(2),
			makeBinaryExpression(
				"<",
				makeIntegerLiteral(2),
				makeBinaryExpression(
					"+",
					makeIntegerLiteral(2),
					makeIntegerLiteral(2),
				),
			),
		)),
		makeExpressionStatement(makeBinaryExpression(
			"+",
			makeUnaryExpression(
				"-",
				makeIntegerLiteral(2),
			),
			makeIntegerLiteral(2),
		)),
	)

	if ast.String() != expectedAst.String() {
		t.Fatalf("Expected: %q, got %q", expectedAst, ast)
	}
}

func TestUnaryExpression(t *testing.T) {
	input := test.MakeInput(
		`-42`,
		`--42`,
		`!eevee`,
		`!!eevee`,
		`!(2 is 2)`,
	)

	l := lexer.New(input, 4)
	p := New(l.Tokens, true)
	ast := p.Parse()

	checkParserErrors(t, p)

	expectedAst := makeProgram(
		makeExpressionStatement(makeUnaryExpression(
			"-",
			makeIntegerLiteral(42),
		)),
		makeExpressionStatement(makeUnaryExpression(
			"-",
			makeUnaryExpression(
				"-",
				makeIntegerLiteral(42),
			),
		)),
		makeExpressionStatement(makeUnaryExpression(
			"!",
			makeIdentifier("eevee"),
		)),
		makeExpressionStatement(makeUnaryExpression(
			"!",
			makeUnaryExpression(
				"!",
				makeIdentifier("eevee"),
			),
		)),
		makeExpressionStatement(makeUnaryExpression(
			"!",
			makeBinaryExpression(
				"==",
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
		`true`,
		`false`,
		`null`,
	)

	l := lexer.New(input, 4)
	p := New(l.Tokens, true)
	ast := p.Parse()

	checkParserErrors(t, p)

	expectedAst := makeProgram(
		makeExpressionStatement(makeIntegerLiteral(42)),
		makeExpressionStatement(makeStringLiteral("eevee")),
		makeExpressionStatement(makeFloatLiteral(3.14)),
		makeExpressionStatement(makeBoolLiteral(true)),
		makeExpressionStatement(makeBoolLiteral(false)),
		makeExpressionStatement(makeNullLiteral()),
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

func makeFunctionDeclaration(name ast.Identifier, params []ast.Identifier, body ast.Statement) *ast.FunctionDeclaration {
	return ast.NewFunctionDeclaration(name, params, body)
}

func makeFunctionParameters(params ...ast.Identifier) []ast.Identifier {
	p := make([]ast.Identifier, 0)

	if len(params) == 0 {
		return p
	}

	p = append(p, params...)

	return p
}

func makeReturnStatement(value ast.Expression) *ast.ReturnStatement {
	return ast.NewReturnStatement(value)
}

func makeVariableStatement(dcls ...*ast.VariableDeclaration) *ast.VariableStatement {
	d := []*ast.VariableDeclaration{}
	d = append(d, dcls...)
	return ast.NewVariableStatement(d)
}

func makeVariableDeclaration(ident ast.Expression, init ast.Expression) *ast.VariableDeclaration {
	return ast.NewVariableDeclaration(ident, init)
}

func makeWhileStatement(cond ast.Expression, body ast.Statement) *ast.WhileStatement {
	return ast.NewWhileStatement(cond, body)
}

func makeDoWhileStatement(cond ast.Expression, body ast.Statement) *ast.DoWhileStatement {
	return ast.NewDoWhileStatement(cond, body)
}

func makeForStatement(init ast.Node, cond, iter ast.Expression, body ast.Statement) *ast.ForStatement {
	return ast.NewForStatement(init, cond, iter, body)
}

func makeIfStatement(cond ast.Expression, cons, alt ast.Statement) *ast.IfStatement {
	return ast.NewIfStatement(cond, cons, alt)
}

func makeExpressionStatement(e ast.Expression) *ast.ExpressionStatement {
	return ast.NewExpressionStatement(e)
}

func makeAssignmentExpression(op string, l, r ast.Expression) *ast.AssignmentExpression {
	return ast.NewAssignmentExpression(op, l, r)
}

func makeLogicalExpression(op string, l, r ast.Expression) *ast.LogicalExpression {
	return ast.NewLogicalExpression(op, l, r)
}

func makeBinaryExpression(op string, l, r ast.Expression) *ast.BinaryExpression {
	return ast.NewBinaryExpression(op, l, r)
}

func makeUnaryExpression(op string, r ast.Expression) *ast.UnaryExpression {
	return ast.NewUnaryExpression(op, r)
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

func makeBoolLiteral(b bool) *ast.BoolLiteral {
	return ast.NewBoolLiteral(b)
}

func makeNullLiteral() *ast.NullLiteral {
	return ast.NewNullLiteral()
}

func makeIdentifier(name string) *ast.Identifier {
	return ast.NewIdentifier(name)
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser error: %q", msg)
	}
	t.FailNow()
}
