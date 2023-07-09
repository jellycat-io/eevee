package parser

import (
	"fmt"
	"strconv"

	"github.com/jellycat-io/eevee/ast"
	"github.com/jellycat-io/eevee/token"
)

var literalTypes = []token.TokenType{
	token.INT,
	token.FLOAT,
	token.STRING,
}

var complexAssignmentOps = []token.TokenType{
	token.PLUS_ASSIGN,
	token.MINUS_ASSIGN,
	token.STAR_ASSIGN,
	token.SLASH_ASSIGN,
	token.PERCENT_ASSIGN,
}

type Parser struct {
	tokens          []token.Token
	currentTokenIdx int
	currentToken    token.Token
	errors          []error
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
		errors:          make([]error, 0),
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
	stmts := p.parseStatements(token.EOF)

	return ast.NewProgram(stmts)
}

func (p *Parser) parseStatements(stopTokens ...token.TokenType) []ast.Statement {
	stmts := []ast.Statement{}

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

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	exp := p.parseExpression()
	return ast.NewExpressionStatement(exp)
}

func (p *Parser) parseExpression() ast.Expression {
	if p.currentToken.Type == token.ILLEGAL {
		t := p.currentToken
		p.advance() // Skip the illegal token
		p.error(fmt.Errorf("[%d:%d] Unexpected token: %q", t.Line, t.Column, t.Type))
		return nil
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
	left := p.parseAdditiveExpression()

	if !isAssignmentOperator(p.currentToken.Type) {
		return left
	}

	return ast.NewAssignmentExpression(
		p.parseAssignmentOperator().Literal,
		left,
		p.parseAssignmentExpression(),
	)
}

func (p *Parser) parseAdditiveExpression() ast.Expression {
	return p.parseBinaryExpression(p.parseMultiplicativeExpression, token.PLUS, token.MINUS)
}

func (p *Parser) parseMultiplicativeExpression() ast.Expression {
	return p.parseBinaryExpression(p.parsePrimaryExpression, token.STAR, token.SLASH, token.PERCENT)
}

func (p *Parser) parseBinaryExpression(builder func() ast.Expression, ops ...token.TokenType) ast.Expression {
	exp := builder()

	for _, op := range ops {
		if p.match(op) {
			operator := p.eat(op)
			right := builder()
			exp = ast.NewBinaryExpression(operator.Literal, exp, right)
			break
		}
	}

	return exp
}

func (p *Parser) parsePrimaryExpression() ast.Expression {
	if isLiteral(p.currentToken.Type) {
		return p.parseLiteral()
	} else if p.match(token.LPAREN) {
		return p.parseGroupedExpression()
	} else if p.match(token.IDENT) {
		return p.parseIdentifier()
	}

	p.error(fmt.Errorf("[%d:%d] Unexpected token: %q", p.currentToken.Line, p.currentToken.Column, p.currentToken.Type))
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
	tok := p.eat(token.IDENT)
	return ast.NewIdentifier(tok.Literal)
}

func (p *Parser) parseLiteral() ast.Expression {
	switch p.currentToken.Type {
	case token.INT:
		return p.parseIntegerLiteral()
	case token.FLOAT:
		return p.parseFloatLiteral()
	case token.STRING:
		return p.parseStringLiteral()
	default:
		p.error(fmt.Errorf("[%d:%d] Unexpected token: %q", p.currentToken.Line, p.currentToken.Column, p.currentToken.Type))
		p.advance()
		return nil
	}
}

func (p *Parser) parseIntegerLiteral() *ast.IntegerLiteral {
	tok := p.eat(token.INT)

	value, err := strconv.ParseInt(tok.Literal, 0, 64)
	if err != nil {
		p.error(fmt.Errorf("[%d:%d] Could not parse %q as integer", p.currentToken.Line, p.currentToken.Column, tok.Literal))
	}

	return ast.NewIntegerLiteral(int64(value))
}

func (p *Parser) parseFloatLiteral() *ast.FloatLiteral {
	tok := p.eat(token.FLOAT)
	value, err := strconv.ParseFloat(tok.Literal, 64)

	if err != nil {
		p.error(fmt.Errorf("[%d:%d] Could not parse %q as float", p.currentToken.Line, p.currentToken.Column, tok.Literal))
	}

	return ast.NewFloatLiteral(float64(value))
}

func (p *Parser) parseStringLiteral() *ast.StringLiteral {
	tok := p.eat(token.STRING)

	return ast.NewStringLiteral(tok.Literal[1 : len(tok.Literal)-1])
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
		p.error(fmt.Errorf("[%d:%d] Expected assignment operator, but got %q", p.currentToken.Line, p.currentToken.Column, p.currentToken.Type))
		return token.ILLEGAL
	}
}

func (p *Parser) eat(tokenType token.TokenType) token.Token {
	tok := p.currentToken
	if !p.match(tokenType) {
		p.error(fmt.Errorf("[%d:%d] Expected %q, but got %q", p.currentToken.Line, p.currentToken.Column, tokenType, p.currentToken.Type))
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

func (p *Parser) error(err error) {
	if p.panicMode {
		return // Don't report multiple errors for the same token
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
	for _, tt := range literalTypes {
		if tt == tokenType {
			return true
		}
	}

	return false
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

	for _, tt := range complexAssignmentOps {
		if tt == tokenType {
			return true
		}
	}

	return false
}
