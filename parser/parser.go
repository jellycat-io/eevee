package parser

import (
	"fmt"
	"strconv"

	"github.com/jellycat-io/eevee/ast"
	"github.com/jellycat-io/eevee/token"
)

var (
	literalTypes = map[token.TokenType]bool{
		token.INT:    true,
		token.FLOAT:  true,
		token.STRING: true,
		token.TRUE:   true,
		token.FALSE:  true,
		token.NULL:   true,
	}
	complexAssignmentOps = map[token.TokenType]bool{
		token.PLUS_ASSIGN:    true,
		token.MINUS_ASSIGN:   true,
		token.STAR_ASSIGN:    true,
		token.SLASH_ASSIGN:   true,
		token.PERCENT_ASSIGN: true,
	}
	nilExpression = ast.NewNullLiteral()
)

type ParseError struct {
	Line    int
	Column  int
	Message string
}

func (pe *ParseError) Error() string {
	return fmt.Errorf("[%d, %d] %s", pe.Line, pe.Column, pe.Message).Error()
}

type Parser struct {
	tokens          []token.Token
	currentTokenIdx int
	currentToken    token.Token
	errors          []ParseError
	panicMode       bool
	isREPL          bool
}

func New(tokens []token.Token, isREPL bool) *Parser {
	currentTokenIdx := 0
	currentToken := tokens[currentTokenIdx]

	return &Parser{
		tokens:          tokens,
		currentTokenIdx: currentTokenIdx,
		currentToken:    currentToken,
		errors:          make([]ParseError, 0),
		isREPL:          isREPL,
	}
}

func (p *Parser) Errors() []string {
	errMsgs := make([]string, len(p.errors))
	for i, err := range p.errors {
		errMsgs[i] = err.Error()
	}

	return errMsgs
}

func (p *Parser) Parse() *ast.Program {
	if len(p.tokens) == 0 {
		return nil
	}

	return p.parseProgram()
}

func (p *Parser) parseProgram() *ast.Program {
	return ast.NewProgram(p.parseStatements(token.EOF))
}

func (p *Parser) parseStatements(stopTokens ...token.TokenType) []ast.Statement {
	stmts := make([]ast.Statement, 0, len(p.tokens))

	for !p.matchAny(stopTokens...) {
		stmts = append(stmts, p.parseStatement())
	}

	return stmts
}

func (p *Parser) parseStatement() ast.Statement {
	var stmt ast.Statement
	switch p.currentToken.Type {
	case token.INDENT:
		stmt = p.parseBlockStatement()
	case token.LET:
		stmt = p.parseVariableStatement()
	case token.IF:
		stmt = p.parseIfStatement()
	case token.WHILE, token.DO, token.FOR:
		stmt = p.parseIterationStatement()
	case token.FUNCTION:
		stmt = p.parseFunctionDeclaration()
	case token.RETURN:
		stmt = p.parseReturnStatement()
	default:
		stmt = p.parseExpressionStatement()
	}

	if p.match(token.EOL) {
		p.eat(token.EOL)
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	stmts := []ast.Statement{}
	p.eat(token.INDENT)

	if !p.match(token.DEDENT) {
		stmts = append(stmts, p.parseStatements(token.DEDENT)...)
	}

	p.eat(token.DEDENT)

	return ast.NewBlockStatement(stmts)
}

func (p *Parser) parseFunctionDeclaration() *ast.FunctionDeclaration {
	p.eat(token.FUNCTION)
	name := *p.parseIdentifier()
	p.eat(token.LPAREN)

	params := p.parseFunctionParameters()
	p.eat(token.RPAREN)

	if p.match(token.EOL) {
		p.eat(token.EOL)
	}

	body := p.parseStatement()

	return ast.NewFunctionDeclaration(name, params, body)
}

func (p *Parser) parseFunctionParameters() []ast.Identifier {
	params := make([]ast.Identifier, 0)

	for !p.match(token.RPAREN) {
		params = append(params, *p.parseIdentifier())
		for p.match(token.COMMA) {
			p.eat(token.COMMA)
			params = append(params, *p.parseIdentifier())
		}
	}

	return params
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	p.eat(token.RETURN)
	if p.match(token.EOL) || p.match(token.DEDENT) || p.isAtEnd() {
		return ast.NewReturnStatement(nilExpression)
	}

	return ast.NewReturnStatement(p.parseExpression())
}

func (p *Parser) parseVariableStatement() *ast.VariableStatement {
	p.eat(token.LET)
	declarations := p.parseVariableDeclarationList()

	return ast.NewVariableStatement(declarations)
}

func (p *Parser) parseVariableDeclarationList() []*ast.VariableDeclaration {
	dcls := []*ast.VariableDeclaration{p.parseVariableDeclaration()}

	for p.match(token.COMMA) {
		p.eat(token.COMMA)
		dcls = append(dcls, p.parseVariableDeclaration())
	}

	return dcls
}

func (p *Parser) parseVariableDeclaration() *ast.VariableDeclaration {
	ident := p.parseIdentifier()
	var init ast.Expression

	if !p.match(token.COMMA) && p.match(token.ASSIGN) {
		init = p.parseVariableInitializer()
	} else {
		init = nil
	}

	return ast.NewVariableDeclaration(ident, init)
}

func (p *Parser) parseVariableInitializer() ast.Expression {
	p.eat(token.ASSIGN)
	return p.parseAssignmentExpression()
}

func (p *Parser) parseIterationStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.WHILE:
		return p.parseWhileStatement()
	case token.DO:
		return p.parseDoWhileStatement()
	case token.FOR:
		return p.parseForStatement()
	default:
		return nil
	}
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	p.eat(token.WHILE)
	cond := p.parseExpression()
	p.eat(token.DO)
	if p.match(token.EOL) {
		p.eat(token.EOL)
	}
	body := p.parseStatement()

	return ast.NewWhileStatement(cond, body)
}

func (p *Parser) parseDoWhileStatement() *ast.DoWhileStatement {
	p.eat(token.DO)
	body := p.parseStatement()
	p.eat(token.WHILE)
	if p.match(token.EOL) {
		p.eat(token.EOL)
	}
	cond := p.parseExpression()

	return ast.NewDoWhileStatement(cond, body)
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	p.eat(token.FOR)

	var init ast.Node
	if !p.match(token.SEMI) {
		init = p.parseForStatementInitializer()
	}
	p.eat(token.SEMI)

	var cond ast.Expression
	if !p.match(token.SEMI) {
		cond = p.parseExpression()
	}
	p.eat(token.SEMI)

	var iter ast.Expression
	if !p.match(token.DO) {
		iter = p.parseExpression()
	}
	p.eat(token.DO)

	if p.match(token.EOL) {
		p.eat(token.EOL)
	}

	body := p.parseStatement()

	return ast.NewForStatement(init, cond, iter, body)
}

func (p *Parser) parseForStatementInitializer() ast.Node {
	if p.match(token.LET) {
		return p.parseVariableStatement()
	}

	return p.parseExpression()
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	p.eat(token.IF)
	condition := p.parseExpression()
	p.eat(token.THEN)

	if p.match(token.EOL) {
		p.eat(token.EOL)
	}

	consequent := p.parseStatement()

	var alternate ast.Statement
	if p.match(token.ELSE) {
		p.eat(token.ELSE)

		if p.match(token.EOL) {
			p.eat(token.EOL)
		}

		alternate = p.parseStatement()
	} else {
		alternate = nil
	}

	return ast.NewIfStatement(condition, consequent, alternate)
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	exp := p.parseExpression()
	return ast.NewExpressionStatement(exp)
}

func (p *Parser) parseExpression() ast.Expression {
	if p.currentToken.Type == token.ILLEGAL {
		t := p.currentToken
		p.advance() // Skip the illegal token
		p.error(t.Line, t.Column, fmt.Sprintf("Unexpected token: %q", t.Type))
		return nilExpression
	}

	var exp ast.Expression

	if p.isREPL && p.peekToken().Type == token.EOF && (isLiteral(p.currentToken.Type) || p.match(token.IDENT)) {
		exp = p.parsePrimaryExpression()
	} else {
		exp = p.parseAssignmentExpression()
	}

	// If we've encountered an error and haven't consumed a new token,
	// enter panic mode and synchronize.
	if p.panicMode {
		p.synchronize()
	}

	return exp
}

func (p *Parser) parseAssignmentExpression() ast.Expression {
	left := p.parseLogicalOrExpression()

	if !isAssignmentOperator(p.currentToken.Type) {
		return left
	}

	return ast.NewAssignmentExpression(
		p.parseAssignmentOperator().Literal,
		left,
		p.parseAssignmentExpression(),
	)
}

func (p *Parser) parseLogicalOrExpression() ast.Expression {
	return p.parseLogicalExpression(p.parseLogicalAndExpression, token.OR)
}

func (p *Parser) parseLogicalAndExpression() ast.Expression {
	return p.parseLogicalExpression(p.parseEqualityExpression, token.AND)
}

func (p *Parser) parseLogicalExpression(builder func() ast.Expression, ops ...token.TokenType) ast.Expression {
	exp := builder()

	for _, op := range ops {
		if p.match(op) {
			op_lit := p.eat(op).Literal
			if op_lit == "and" {
				op_lit = "&&"
			}
			if op_lit == "or" {
				op_lit = "||"
			}
			right := builder()
			exp = ast.NewLogicalExpression(op_lit, exp, right)
			break
		}
	}

	return exp
}

func (p *Parser) parseEqualityExpression() ast.Expression {
	return p.parseBinaryExpression(p.parseRelationalExpression, token.EQ, token.NOT_EQ)
}

func (p *Parser) parseRelationalExpression() ast.Expression {
	return p.parseBinaryExpression(p.parseAdditiveExpression, token.LT, token.LT_EQ, token.GT, token.GT_EQ)
}

func (p *Parser) parseAdditiveExpression() ast.Expression {
	return p.parseBinaryExpression(p.parseMultiplicativeExpression, token.PLUS, token.MINUS)
}

func (p *Parser) parseMultiplicativeExpression() ast.Expression {
	return p.parseBinaryExpression(p.parseUnaryExpression, token.STAR, token.SLASH, token.PERCENT)
}

func (p *Parser) parseBinaryExpression(builder func() ast.Expression, ops ...token.TokenType) ast.Expression {
	exp := builder()

	for _, op := range ops {
		if p.match(op) {
			op_lit := p.eat(op).Literal
			if op_lit == "is" {
				op_lit = "=="
			}
			if op_lit == "not" {
				op_lit = "!="
			}
			right := builder()
			exp = ast.NewBinaryExpression(op_lit, exp, right)
			break
		}
	}

	return exp
}

func (p *Parser) parseUnaryExpression() ast.Expression {
	var op string
	switch p.currentToken.Type {
	case token.PLUS:
		op = p.eat(token.PLUS).Literal
	case token.MINUS:
		op = p.eat(token.MINUS).Literal
	case token.BANG:
		op = p.eat(token.BANG).Literal
	}

	if op != "" {
		return ast.NewUnaryExpression(op, p.parseUnaryExpression())
	}

	return p.parseLeftHandSideExpression()
}

func (p *Parser) parseLeftHandSideExpression() ast.Expression {
	return p.parsePrimaryExpression()
}

func (p *Parser) parsePrimaryExpression() ast.Expression {
	if isLiteral(p.currentToken.Type) {
		return p.parseLiteral()
	} else if p.match(token.LPAREN) {
		return p.parseGroupedExpression()
	} else if p.match(token.IDENT) {
		return p.parseIdentifier()
	}

	p.error(p.currentToken.Line, p.currentToken.Column, fmt.Sprintf("Unexpected token: %q", p.currentToken.Type))
	p.advance()
	return nil
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.eat(token.LPAREN)
	exp := p.parseExpression()
	p.eat(token.RPAREN)

	return exp
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	return ast.NewIdentifier(p.eat(token.IDENT).Literal)
}

func (p *Parser) parseLiteral() ast.Expression {
	switch p.currentToken.Type {
	case token.INT:
		return p.parseIntegerLiteral()
	case token.FLOAT:
		return p.parseFloatLiteral()
	case token.STRING:
		return p.parseStringLiteral()
	case token.TRUE:
		return p.parseBoolLiteral(true)
	case token.FALSE:
		return p.parseBoolLiteral(false)
	case token.NULL:
		return p.parseNullLiteral()
	default:
		p.error(p.currentToken.Line, p.currentToken.Column, fmt.Sprintf("Unexpected token: %q", p.currentToken.Type))
		p.advance()
		return nilExpression
	}
}

func (p *Parser) parseIntegerLiteral() *ast.IntegerLiteral {
	tok := p.eat(token.INT)

	value, err := strconv.ParseInt(tok.Literal, 0, 64)
	if err != nil {
		p.error(p.currentToken.Line, p.currentToken.Column, fmt.Sprintf("Could not parse %q as integer", tok.Literal))
	}

	return ast.NewIntegerLiteral(int64(value))
}

func (p *Parser) parseFloatLiteral() *ast.FloatLiteral {
	tok := p.eat(token.FLOAT)
	value, err := strconv.ParseFloat(tok.Literal, 64)

	if err != nil {
		p.error(p.currentToken.Line, p.currentToken.Column, fmt.Sprintf("Could not parse %q as float", tok.Literal))
	}

	return ast.NewFloatLiteral(float64(value))
}

func (p *Parser) parseStringLiteral() *ast.StringLiteral {
	tok := p.eat(token.STRING)

	return ast.NewStringLiteral(tok.Literal[1 : len(tok.Literal)-1])
}

func (p *Parser) parseBoolLiteral(value bool) *ast.BoolLiteral {
	switch value {
	case true:
		p.eat(token.TRUE)
	case false:
		p.eat(token.FALSE)
	}

	return ast.NewBoolLiteral(value)
}

func (p *Parser) parseNullLiteral() *ast.NullLiteral {
	p.eat(token.NULL)
	return nilExpression
}

func (p *Parser) parseAssignmentOperator() token.Token {
	if p.match(token.ASSIGN) {
		tok := p.eat(token.ASSIGN)
		return tok
	}

	tokenType := p.checkComplexAssignmentOperator()
	tok := p.eat(tokenType)
	return tok
}

func (p *Parser) checkComplexAssignmentOperator() token.TokenType {
	switch p.currentToken.Type {
	case token.PLUS_ASSIGN:
		return token.PLUS_ASSIGN
	case token.MINUS_ASSIGN:
		return token.MINUS_ASSIGN
	case token.STAR_ASSIGN:
		return token.STAR_ASSIGN
	case token.SLASH_ASSIGN:
		return token.SLASH_ASSIGN
	case token.PERCENT_ASSIGN:
		return token.PERCENT_ASSIGN
	default:
		p.error(p.currentToken.Line, p.currentToken.Column, fmt.Sprintf("Expected assignment operator, but got %q", p.currentToken.Type))
		return token.ILLEGAL
	}
}

func (p *Parser) eat(tokenType token.TokenType) token.Token {
	tok := p.currentToken
	if !p.match(tokenType) {
		p.error(p.currentToken.Line, p.currentToken.Column, fmt.Sprintf("Expected %q, but got %q", tokenType, p.currentToken.Type))
		tok = token.NewToken(token.ILLEGAL, "", p.currentToken.Line, p.currentToken.Column)
	}

	p.advance()

	return tok
}

func (p *Parser) synchronize() {
	for !p.match(token.EOL) && !p.isAtEnd() {
		p.advance()
	}

	if !p.isREPL && p.isAtEnd() {
		p.eat(token.EOL)
	}
}

func (p *Parser) advance() {
	p.currentTokenIdx++

	if p.currentTokenIdx < len(p.tokens) {
		p.currentToken = p.tokens[p.currentTokenIdx]
	} else {
		p.currentToken = token.NewToken(token.EOF, "", p.currentToken.Line, p.currentToken.Column+1)
	}

	p.panicMode = false
}

func (p *Parser) error(line, column int, msg string) {
	if p.panicMode {
		return // Don't report multiple errors for the same token
	}
	err := ParseError{
		Line:    line,
		Column:  column,
		Message: msg,
	}
	p.errors = append(p.errors, err)

	p.panicMode = true // Enter panic mode after an error
}

func (p *Parser) match(tokenType token.TokenType) bool {
	return p.currentToken.Type == tokenType
}

func (p *Parser) matchAny(tokenTypes ...token.TokenType) bool {
	for _, tt := range tokenTypes {
		if p.match(tt) {
			return true
		}
	}
	return false
}

func (p *Parser) peekToken() token.Token {
	if p.currentTokenIdx+1 < len(p.tokens) {
		return p.tokens[p.currentTokenIdx+1]
	}
	return token.Token{}
}

func (p *Parser) isAtEnd() bool {
	return p.currentToken.Type == token.EOF
}

func isLiteral(tokenType token.TokenType) bool {
	return literalTypes[tokenType]
}

// func isBinaryOperator(tokenType token.TokenType) bool {
// 	for _, tt := range binaryOperators {
// 		if tt == tokenType {
// 			return true
// 		}
// 	}

// 	return false
// }

func isAssignmentOperator(tokenType token.TokenType) bool {
	if tokenType == token.ASSIGN {
		return true
	}

	if complexAssignmentOps[tokenType] {
		return true
	}

	return false
}
