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

type ExpressionStatement struct {
	Expression Expression `json:"expression"`
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) String() string {
	return fmt.Sprintf("ExpressionStatement(%v)", es.Expression)
}

func NewExpressionStatement(exp Expression) *ExpressionStatement {
	return &ExpressionStatement{Expression: exp}
}

type IntegerLiteral struct {
	Value int64 `json:"value"`
}

func (il *IntegerLiteral) expressionNode() {}
func (il IntegerLiteral) String() string {
	return fmt.Sprintf("IntegerLiteral(%d)", il.Value)
}

func NewIntegerLiteral(value int64) *IntegerLiteral {
	return &IntegerLiteral{Value: value}
}

type FloatLiteral struct {
	Value float64 `json:"value"`
}

func (fl *FloatLiteral) expressionNode() {}
func (fl FloatLiteral) String() string {
	return fmt.Sprintf("FloatLiteral(%v)", fl.Value)
}

func NewFloatLiteral(value float64) *FloatLiteral {
	return &FloatLiteral{Value: value}
}

type StringLiteral struct {
	Value string `json:"value"`
}

func (sl *StringLiteral) expressionNode() {}
func (sl StringLiteral) String() string {
	return fmt.Sprintf("StringLiteral(%s)", sl.Value)
}

func NewStringLiteral(value string) *StringLiteral {
	return &StringLiteral{Value: value}
}
