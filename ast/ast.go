package ast

import (
	"fmt"
)

type Node interface {
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	// program ::= statements EOF
	Type       string      `json:"type"`
	Statements []Statement `json:"statements"`
}

func (p *Program) String() string {
	result := "(Program("
	for _, stmt := range p.Statements {
		result += stmt.String()
	}
	result += "))"

	return result
}

func NewProgram(statements []Statement) *Program {
	return &Program{
		Type:       "Program",
		Statements: statements,
	}
}

type BlockStatement struct {
	// block_statement ::= INDENT statements DEDENT
	Type       string      `json:"type"`
	Statements []Statement `json:"statements"`
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) String() string {
	result := "BlockStatement("
	for _, stmt := range bs.Statements {
		result += stmt.String()
	}
	result += ")"

	return result
}

func NewBlockStatement(statements []Statement) *BlockStatement {
	return &BlockStatement{
		Type:       "BlockStatement",
		Statements: statements,
	}
}

type ExpressionStatement struct {
	// expression_statement ::= expression
	Type       string     `json:"type"`
	Expression Expression `json:"expression"`
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) String() string {
	return fmt.Sprintf("ExpressionStatement(%v)", es.Expression)
}

func NewExpressionStatement(exp Expression) *ExpressionStatement {
	return &ExpressionStatement{Type: "ExpressionStatement", Expression: exp}
}

type BinaryExpression struct {
	// additive_expression 			::= multiplicative_expression { (PLUS | MINUS) multiplicative_expression }
	// multiplicative_expression 	::= primary_expression { (STAR | SLASH | PERCENT) primary_expression }
	Type     string     `json:"type"`
	Operator string     `json:"operator"`
	Left     Expression `json:"left"`
	Right    Expression `json:"right"`
}

func (be *BinaryExpression) expressionNode() {}
func (be *BinaryExpression) String() string {
	return fmt.Sprintf("BinaryExpression(%s %v %v)", be.Operator, be.Left, be.Right)
}

func NewBinaryExpression(op string, left, right Expression) *BinaryExpression {
	return &BinaryExpression{
		Type:     "BinaryExpression",
		Operator: op,
		Left:     left,
		Right:    right,
	}
}

type IntegerLiteral struct {
	// integer_literal ::= INT
	Type  string `json:"type"`
	Value int64  `json:"value"`
}

func (il *IntegerLiteral) expressionNode() {}
func (il IntegerLiteral) String() string {
	return fmt.Sprintf("IntegerLiteral(%d)", il.Value)
}

func NewIntegerLiteral(value int64) *IntegerLiteral {
	return &IntegerLiteral{Type: "IntegerLiteral", Value: value}
}

type FloatLiteral struct {
	// float_literal ::= FLOAT
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

func (fl *FloatLiteral) expressionNode() {}
func (fl FloatLiteral) String() string {
	return fmt.Sprintf("FloatLiteral(%v)", fl.Value)
}

func NewFloatLiteral(value float64) *FloatLiteral {
	return &FloatLiteral{Type: "FloatLiteral", Value: value}
}

type StringLiteral struct {
	// string_literal ::= STRING
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (sl *StringLiteral) expressionNode() {}
func (sl StringLiteral) String() string {
	return fmt.Sprintf("StringLiteral(%s)", sl.Value)
}

func NewStringLiteral(value string) *StringLiteral {
	return &StringLiteral{Type: "StringLiteral", Value: value}
}
