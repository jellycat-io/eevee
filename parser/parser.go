package parser

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jellycat-io/eevee/ast"
	"github.com/jellycat-io/eevee/token"
)

type Parser struct {
	tokens          []token.Token
	currentTokenIdx int
	currentToken    token.Token
}

func NewParser(tokens []token.Token) *Parser {
	currentTokenIdx := 0
	currentToken := tokens[currentTokenIdx]

	return &Parser{
		tokens:          tokens,
		currentTokenIdx: currentTokenIdx,
		currentToken:    currentToken,
	}
}

func (p *Parser) Parse() *ast.Program {
	if len(p.tokens) == 0 {
		return &ast.Program{}
	}

	return p.parseProgram()
}

func (p *Parser) parseProgram() *ast.Program {
	stmts := p.parseStatements(token.EOF)

	return ast.NewProgram(stmts)
}

func (p *Parser) parseStatements(stopTokenType token.TokenType) []ast.Statement {
	stmts := []ast.Statement{}

	for !p.match(stopTokenType) {
		stmts = append(stmts, p.parseStatement())
	}

	return stmts
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.INDENT:
		return p.parseBlockStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	stmts := []ast.Statement{}
	if _, err := p.eat(token.INDENT); err != nil {
		log.Fatal(err)
	}

	if !p.match(token.DEDENT) {
		stmts = append(stmts, p.parseStatements(token.DEDENT)...)
	}

	if _, err := p.eat(token.DEDENT); err != nil {
		log.Fatal(err)
	}

	return ast.NewBlockStatement(stmts)
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	exp := p.parseExpression()
	return ast.NewExpressionStatement(exp)
}

func (p *Parser) parseExpression() ast.Expression {
	return p.parsePrimaryExpression()
}

func (p *Parser) parsePrimaryExpression() ast.Expression {
	lit, err := p.parseLiteral()
	if err != nil {
		log.Fatal(err)
	}

	return lit
}

func (p *Parser) parseLiteral() (ast.Expression, error) {
	switch p.currentToken.Type {
	case token.INT:
		return p.parseIntegerLiteral(), nil
	case token.FLOAT:
		return p.parseFloatLiteral(), nil
	case token.STRING:
		return p.parseStringLiteral(), nil
	default:
		return nil, fmt.Errorf("[%d:%d] Unexpected token: %q", p.currentToken.Line, p.currentToken.Column, p.currentToken.Type)
	}
}

func (p *Parser) parseIntegerLiteral() *ast.IntegerLiteral {
	tok, err := p.eat(token.INT)
	if err != nil {
		log.Fatal(err)
	}
	value, err := strconv.ParseInt(tok.Literal, 0, 64)
	if err != nil {
		log.Fatalf("[%d:%d] Could not parse %q as integer", p.currentToken.Line, p.currentToken.Column, tok.Literal)
	}

	return ast.NewIntegerLiteral(int64(value))
}

func (p *Parser) parseFloatLiteral() *ast.FloatLiteral {
	tok, err := p.eat(token.FLOAT)
	if err != nil {
		log.Fatal(err)
	}
	value, err := strconv.ParseFloat(tok.Literal, 64)
	if err != nil {
		log.Fatalf("[%d:%d] Could not parse %q as float", p.currentToken.Line, p.currentToken.Column, tok.Literal)
	}

	return ast.NewFloatLiteral(float64(value))
}

func (p *Parser) parseStringLiteral() *ast.StringLiteral {
	tok, err := p.eat(token.STRING)
	if err != nil {
		log.Fatal(err)
	}

	return ast.NewStringLiteral(tok.Literal[1 : len(tok.Literal)-1])
}

func (p *Parser) eat(tokenType token.TokenType) (token.Token, error) {
	tok := p.currentToken

	if p.match(tokenType) {
		p.advance()
	} else {
		return token.Token{}, fmt.Errorf("[%d:%d] Expected %q, but got %q", p.currentToken.Line, p.currentToken.Column, tokenType, p.currentToken.Type)
	}

	return tok, nil
}

func (p *Parser) advance() {
	p.currentTokenIdx++

	if p.currentTokenIdx < len(p.tokens) {
		p.currentToken = p.tokens[p.currentTokenIdx]
	} else {
		p.currentToken = token.Token{}
	}
}

func (p *Parser) match(tokenType token.TokenType) bool {
	return p.currentToken.Type == tokenType
}

// func (p *Parser) isAtEnd() bool {
// 	return p.currentToken == token.Token{} || p.currentToken.Type == token.EOF
// }
