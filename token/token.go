package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

func NewToken(tokType TokenType, literal string, line, column int) Token {
	return Token{
		Type:    tokType,
		Literal: literal,
		Line:    line,
		Column:  column,
	}
}

const (
	ILLEGAL        = TokenType("ILLEGAL")
	WHITESPACE     = TokenType("WHITESPACE")
	EOL            = TokenType("EOL")
	COMMENT        = TokenType("COMMENT")
	INDENT         = TokenType("INDENT")
	DEDENT         = TokenType("DEDENT")
	EOF            = TokenType("EOF")
	IDENT          = TokenType("IDENT")
	ASSIGN         = TokenType("=")
	AND            = TokenType("&&")
	OR             = TokenType("||")
	EQ             = TokenType("==")
	NOT_EQ         = TokenType("!=")
	LT             = TokenType("<")
	LT_EQ          = TokenType("<=")
	GT             = TokenType(">")
	GT_EQ          = TokenType(">=")
	PLUS_ASSIGN    = TokenType("+=")
	MINUS_ASSIGN   = TokenType("-=")
	STAR_ASSIGN    = TokenType("*=")
	SLASH_ASSIGN   = TokenType("/=")
	PERCENT_ASSIGN = TokenType("%=")
	PLUS           = TokenType("+")
	MINUS          = TokenType("-")
	STAR           = TokenType("*")
	SLASH          = TokenType("/")
	PERCENT        = TokenType("%")
	BANG           = TokenType("!")
	COMMA          = TokenType(",")
	DOT            = TokenType(".")
	SEMI           = TokenType(";")
	COLON          = TokenType(":")
	LPAREN         = TokenType("(")
	RPAREN         = TokenType(")")
	LBRACE         = TokenType("{")
	RBRACE         = TokenType("}")
	LBRACKET       = TokenType("[")
	RBRACKET       = TokenType("]")
	FUNCTION       = TokenType("FUNCTION")
	MODULE         = TokenType("MODULE")
	IMPORT         = TokenType("IMPORT")
	LET            = TokenType("LET")
	TRUE           = TokenType("TRUE")
	FALSE          = TokenType("FALSE")
	IF             = TokenType("IF")
	THEN           = TokenType("THEN")
	ELSE           = TokenType("ELSE")
	WHILE          = TokenType("WHILE")
	FOR            = TokenType("FOR")
	DO             = TokenType("DO")
	RETURN         = TokenType("RETURN")
	NULL           = TokenType("NULL")
	INT            = TokenType("INT")
	FLOAT          = TokenType("FLOAT")
	STRING         = TokenType("STRING")
)

var Keywords = map[string]TokenType{
	"and":    AND,
	"do":     DO,
	"else":   ELSE,
	"false":  FALSE,
	"fn":     FUNCTION,
	"for":    FOR,
	"if":     IF,
	"is":     EQ,
	"import": IMPORT,
	"let":    LET,
	"module": MODULE,
	"null":   NULL,
	"not":    NOT_EQ,
	"or":     OR,
	"return": RETURN,
	"then":   THEN,
	"true":   TRUE,
	"while":  WHILE,
}
