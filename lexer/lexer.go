package lexer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/TwiN/go-color"
	"github.com/jellycat-io/eevee/token"
)

type Lexer struct {
	source      string
	tabSize     int
	Tokens      []token.Token
	indentStack []int
	patterns    []struct {
		regex   *regexp.Regexp
		tokType token.TokenType
	}
}

func New(source string, tabSize int) *Lexer {
	patternList := []struct {
		pattern string
		tokType token.TokenType
	}{
		/****************************************
		 * Whitespaces
		 ***************************************/
		{`^[ \t]+`, token.WHITESPACE},
		/****************************************
		 * Comments
		 ***************************************/
		{`^#[^\n]*`, token.COMMENT},
		/****************************************
		 * Logical operators
		 ***************************************/
		{`^&&`, token.AND},
		{`^\|\|`, token.OR},
		/****************************************
		 * Comparison operators
		 ***************************************/
		{`^==`, token.EQ},
		{`^!=`, token.NOT_EQ},
		{`^<=`, token.LT_EQ},
		{`^>=`, token.GT_EQ},
		{`^<`, token.LT},
		{`^>`, token.GT},
		/****************************************
		 * Symbols, delimiters
		 ***************************************/
		{`^;`, token.SEMI},
		{`^,`, token.COMMA},
		{`^:`, token.COLON},
		{`^\(`, token.LPAREN},
		{`^\)`, token.RPAREN},
		{`^{`, token.LBRACE},
		{`^}`, token.RBRACE},
		{`^\[`, token.LBRACKET},
		{`^]`, token.RBRACKET},
		{`^!`, token.BANG},
		/****************************************
		 * Identifiers
		 ***************************************/
		{`^[a-zA-Z_][a-zA-Z0-9_]*`, token.IDENT},
		/****************************************
		 * Assignment operators
		 ***************************************/
		{`^=`, token.ASSIGN},
		{`^\+=`, token.PLUS_ASSIGN},
		{`^-=`, token.MINUS_ASSIGN},
		{`^\*=`, token.STAR_ASSIGN},
		{`^/=`, token.SLASH_ASSIGN},
		{`^%=`, token.PERCENT_ASSIGN},
		/****************************************
		 * Math operators
		 ***************************************/
		{`^\+`, token.PLUS},
		{`^-`, token.MINUS},
		{`^\*`, token.STAR},
		{`^/`, token.SLASH},
		{`^%`, token.PERCENT},
		/****************************************
		 * Literals
		 ***************************************/
		{`^\d+\.\d+`, token.FLOAT},
		{`^\d+`, token.INT},
		{`^\".*?\"`, token.STRING},
	}

	compiledPatterns := make([]struct {
		regex   *regexp.Regexp
		tokType token.TokenType
	}, len(patternList))

	for i, pat := range patternList {
		compiled, err := regexp.Compile(pat.pattern)
		if err != nil {
			fmt.Printf(color.InRed("Failed to compile regex pattern %q: %v\n"), pat.pattern, err)
			continue
		}
		compiledPatterns[i].regex = compiled
		compiledPatterns[i].tokType = pat.tokType
	}

	l := &Lexer{
		source:      source,
		indentStack: []int{0},
		tabSize:     tabSize,
		patterns:    compiledPatterns,
	}

	l.tokenize()

	return l
}

func (l *Lexer) tokenize() {
	lines := strings.Split(l.source, "\n")
	for lineNum, line := range lines {
		lineNum++
		column := 1
		indentLevel := len(line) - len(strings.TrimSpace(line))
		indentString := line[:indentLevel]

		if indentLevel == l.indentStack[len(l.indentStack)-1] {
			column += len(indentString)
		}

		for indentLevel < l.indentStack[len(l.indentStack)-1] {
			l.Tokens = append(l.Tokens, token.NewToken(token.DEDENT, "", lineNum, column))
			l.indentStack = l.indentStack[:len(l.indentStack)-1]
		}

		if indentLevel > l.indentStack[len(l.indentStack)-1] {
			l.Tokens = append(l.Tokens, token.NewToken(token.INDENT, "", lineNum, column))
			l.indentStack = append(l.indentStack, indentLevel)
			for _, char := range indentString {
				if char == '\t' {
					spacesNeeded := l.tabSize - ((column - 1) % l.tabSize)
					column += spacesNeeded
				} else if char == ' ' {
					column++
				}
			}
		}

		column = l.tokenizeLine(line, lineNum, column)

		if lineNum != len(lines) { // Skip adding EOL token for the last line
			l.Tokens = append(l.Tokens, token.NewToken(token.EOL, "", lineNum, column))
		}
	}

	for range l.indentStack[1:] {
		l.Tokens = append(l.Tokens, token.NewToken(token.DEDENT, "", len(lines)+1, 1))
	}

	l.Tokens = append(l.Tokens, token.NewToken(token.EOF, "", len(lines)+1, 1))
}

func (l *Lexer) tokenizeLine(line string, lineNum, column int) int {
	line = strings.TrimSpace(line)

	for line != "" {
		matched := false
		for _, pattern := range l.patterns {
			loc := pattern.regex.FindStringIndex(line)
			if loc != nil {
				lexeme := line[loc[0]:loc[1]]

				tokenType := pattern.tokType
				if tokenType == token.IDENT {
					tokenType = lookupIdent(lexeme)
				}

				if tokenType != token.WHITESPACE && tokenType != token.COMMENT {
					l.Tokens = append(l.Tokens, token.NewToken(tokenType, lexeme, lineNum, column))
				}
				line = line[loc[1]:]
				column += loc[1]
				matched = true
				break
			}
		}

		if !matched {
			l.Tokens = append(l.Tokens, token.NewToken(token.ILLEGAL, line[0:1], lineNum, column))
			line = line[1:]
			column++
		}
	}

	return column
}

func lookupIdent(ident string) token.TokenType {
	if tok, ok := token.Keywords[ident]; ok {
		return tok
	}

	return token.IDENT
}
