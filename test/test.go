package test

import (
	"bytes"
	"strings"

	"github.com/jellycat-io/eevee/token"
)

// TokensFromString tokenizes an input string and returns a slice of Tokens.
func TokensFromString(input string) []token.Token {
	tokenMap := map[string]token.TokenType{
		"INDENT": token.INDENT,
		"DEDENT": token.DEDENT,
		"EOF":    token.EOF,
		"IDENT":  token.IDENT,
		"INT":    token.INT,
		"FLOAT":  token.FLOAT,
		"STRING": token.STRING,
		"=":      token.ASSIGN,
		"+=":     token.PLUS_ASSIGN,
		"-=":     token.MINUS_ASSIGN,
		"*=":     token.STAR_ASSIGN,
		"/=":     token.SLASH_ASSIGN,
		"+":      token.PLUS,
		"-":      token.MINUS,
		"*":      token.STAR,
		"/":      token.SLASH,
		"%":      token.PERCENT,
		"!":      token.BANG,
		"==":     token.EQ,
		"!=":     token.NOT_EQ,
		"<":      token.LT,
		"<=":     token.LT_EQ,
		">":      token.GT,
		">=":     token.GT_EQ,
		"&&":     token.AND,
		"||":     token.OR,
		",":      token.COMMA,
		";":      token.SEMI,
		":":      token.COLON,
		"(":      token.LPAREN,
		")":      token.RPAREN,
		"{":      token.LBRACE,
		"}":      token.RBRACE,
		"[":      token.LBRACKET,
		"]":      token.RBRACKET,
		"fn":     token.FUNCTION,
		"module": token.MODULE,
		"import": token.IMPORT,
		"let":    token.LET,
		"true":   token.TRUE,
		"false":  token.FALSE,
		"if":     token.IF,
		"is":     token.EQ,
		"not":    token.NOT_EQ,
		"and":    token.AND,
		"or":     token.OR,
		"then":   token.THEN,
		"else":   token.ELSE,
		"return": token.RETURN,
		"nil":    token.NIL,
	}

	var tokens []token.Token
	lines := strings.Split(input, "\n")
	currentLine := 1
	indentLevels := []int{0}

	for _, line := range lines {
		indentation := len(line) - len(strings.TrimSpace(line))
		currentIndentLevel := indentLevels[len(indentLevels)-1]

		if indentation > currentIndentLevel {
			tokens = append(tokens, token.Token{Type: token.INDENT, Literal: "", Line: currentLine, Column: 1})
			indentLevels = append(indentLevels, indentation)
		}

		for indentation < currentIndentLevel {
			tokens = append(tokens, token.Token{Type: token.DEDENT, Literal: "", Line: currentLine, Column: 1})
			indentLevels = indentLevels[:len(indentLevels)-1]
			currentIndentLevel = indentLevels[len(indentLevels)-1]
		}

		currentColumn := indentation + 1
		segments := strings.Fields(strings.TrimSpace(line))
		for _, segment := range segments {
			tokenType, literal := determineTokenTypeAndLiteral(tokenMap, segment)

			tokens = append(tokens, token.Token{Type: tokenType, Literal: literal, Line: currentLine, Column: currentColumn})
			currentColumn += len(segment) + 1
		}
		currentLine++
	}

	for len(indentLevels) > 1 {
		tokens = append(tokens, token.Token{Type: token.DEDENT, Literal: "", Line: currentLine, Column: 1})
		indentLevels = indentLevels[:len(indentLevels)-1]
	}

	tokens = append(tokens, token.Token{Type: token.EOF, Literal: "", Line: currentLine, Column: 1})
	return tokens
}

// MakeInput takes code lines as strings or returns them in a single string separated by '\n'
func MakeInput(lines ...string) string {
	var out bytes.Buffer
	for _, line := range lines {
		out.WriteString(line)
		out.WriteString("\n")
	}
	return out.String()
}

func determineTokenTypeAndLiteral(tokenMap map[string]token.TokenType, segment string) (token.TokenType, string) {
	if strings.Contains(segment, "=") && len(segment) > 2 {
		parts := strings.Split(segment, "=")
		typeStr := parts[0]
		literal := parts[1]
		return tokenMap[typeStr], literal
	}
	return tokenMap[segment], segment
}
