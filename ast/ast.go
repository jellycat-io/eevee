package ast

import (
	"fmt"
	"strings"
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
	var result strings.Builder
	result.WriteString("(Program ")
	for _, stmt := range p.Statements {
		result.WriteString(stmt.String())
		result.WriteString(" ")
	}
	result.WriteString(")")
	return strings.TrimSpace(result.String())
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
	var result strings.Builder
	result.WriteString("(BlockStatement ")
	for _, stmt := range bs.Statements {
		result.WriteString(stmt.String())
		result.WriteString(" ")
	}
	result.WriteString(")")
	return strings.TrimSpace(result.String())
}

func NewBlockStatement(statements []Statement) *BlockStatement {
	return &BlockStatement{
		Type:       "BlockStatement",
		Statements: statements,
	}
}

type FunctionDeclaration struct {
	// function_declaration ::= FUNCTION identifier LPAREN parameters RPAREN statement
	// parameters           ::= identifier { COMMA identifier }
	Type       string       `json:"type"`
	Name       Identifier   `json:"name"`
	Parameters []Identifier `json:"parameters"`
	Body       Statement    `json:"body"`
}

func (fd *FunctionDeclaration) statementNode() {}
func (fd *FunctionDeclaration) String() string {
	var result strings.Builder
	result.WriteString("(FunctionDeclaration ")
	result.WriteString(fd.Name.String() + " ")
	for _, param := range fd.Parameters {
		result.WriteString(param.String())
		result.WriteString(" ")
	}
	result.WriteString(fd.Body.String())
	result.WriteString(")")
	return strings.TrimSpace(result.String())
}

func NewFunctionDeclaration(name Identifier, parameters []Identifier, body Statement) *FunctionDeclaration {
	return &FunctionDeclaration{
		Type:       "FunctionDeclaration",
		Name:       name,
		Parameters: parameters,
		Body:       body,
	}
}

type ReturnStatement struct {
	Type  string     `json:"type"`
	Value Expression `json:"value"`
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) String() string {
	return fmt.Sprintf("(ReturnStatement %v)", rs.Value)
}

func NewReturnStatement(value Expression) *ReturnStatement {
	return &ReturnStatement{
		Type:  "ReturnStatement",
		Value: value,
	}
}

type VariableStatement struct {
	// variable_statement ::= LET variable_declaration_list
	// variable_declaration_list ::= variable_declaration { COMMA variable_declaration }
	Type         string                 `json:"type"`
	Declarations []*VariableDeclaration `json:"declarations"`
}

func (vs *VariableStatement) statementNode() {}
func (vs *VariableStatement) String() string {
	var result strings.Builder
	result.WriteString("(VariableStatement ")
	for _, decl := range vs.Declarations {
		result.WriteString(decl.String())
		result.WriteString(" ")
	}
	result.WriteString(")")
	return strings.TrimSpace(result.String())
}

func NewVariableStatement(declarations []*VariableDeclaration) *VariableStatement {
	return &VariableStatement{
		Type:         "VariableStatement",
		Declarations: declarations,
	}
}

type VariableDeclaration struct {
	// variable_declaration ::= identifier [ ASSIGN assignment_expression ]
	Type        string     `json:"type"`
	Identifier  Expression `json:"identifier"`
	Initializer Expression `json:"initializer"`
}

func (vd *VariableDeclaration) String() string {
	return fmt.Sprintf("(VariableDeclaration %v %v)", vd.Identifier, vd.Initializer)
}

func NewVariableDeclaration(identifier Expression, initializer Expression) *VariableDeclaration {
	return &VariableDeclaration{
		Type:        "VariableDeclaration",
		Identifier:  identifier,
		Initializer: initializer,
	}
}

type IfStatement struct {
	// if_statement ::= IF expression THEN statement [ ELSE statement ]
	Type       string     `json:"type"`
	Condition  Expression `json:"condition"`
	Consequent Statement  `json:"consequent"`
	Alternate  Statement  `json:"alternate"`
}

func (is *IfStatement) statementNode() {}
func (is *IfStatement) String() string {
	return fmt.Sprintf("(IfStatement %v %v %v)", is.Condition, is.Consequent, is.Alternate)
}
func NewIfStatement(condition Expression, consequent, alternate Statement) *IfStatement {
	return &IfStatement{
		Type:       "IfStatement",
		Condition:  condition,
		Consequent: consequent,
		Alternate:  alternate,
	}
}

type WhileStatement struct {
	// while_statement ::= WHILE expression DO statement
	Type      string     `json:"type"`
	Condition Expression `json:"condition"`
	Body      Statement  `json:"body"`
}

func (ws *WhileStatement) statementNode() {}
func (ws *WhileStatement) String() string {
	return fmt.Sprintf("(WhileStatement %v %v)", ws.Condition, ws.Body)
}
func NewWhileStatement(condition Expression, body Statement) *WhileStatement {
	return &WhileStatement{
		Type:      "WhileStatement",
		Condition: condition,
		Body:      body,
	}
}

type DoWhileStatement struct {
	// do_while_statement ::= DO statement WHILE expression
	Type      string     `json:"type"`
	Condition Expression `json:"condition"`
	Body      Statement  `json:"body"`
}

func (dws *DoWhileStatement) statementNode() {}
func (dws *DoWhileStatement) String() string {
	return fmt.Sprintf("(DoWhileStatement %v %v)", dws.Condition, dws.Body)
}
func NewDoWhileStatement(condition Expression, body Statement) *DoWhileStatement {
	return &DoWhileStatement{
		Type:      "DoWhileStatement",
		Condition: condition,
		Body:      body,
	}
}

type ForStatement struct {
	// for_statement ::= FOR [ for_statement_initializer ] SEMI [ expression ] SEMI [ expression ] DO statement
	Type        string     `json:"type"`
	Initializer Node       `json:"initializer"`
	Condition   Expression `json:"condition"`
	Iterator    Expression `json:"iterator"`
	Body        Statement  `json:"body"`
}

func (fs *ForStatement) statementNode() {}
func (fs *ForStatement) String() string {
	return fmt.Sprintf("(ForStatement %v %v %v %v)", fs.Initializer, fs.Condition, fs.Iterator, fs.Body)
}
func NewForStatement(initializer Node, condition, iterator Expression, body Statement) *ForStatement {
	return &ForStatement{
		Type:        "ForStatement",
		Initializer: initializer,
		Condition:   condition,
		Iterator:    iterator,
		Body:        body,
	}
}

type ExpressionStatement struct {
	// expression_statement ::= expression
	Type       string     `json:"type"`
	Expression Expression `json:"expression"`
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) String() string {
	return fmt.Sprintf("(ExpressionStatement %v)", es.Expression)
}

func NewExpressionStatement(exp Expression) *ExpressionStatement {
	return &ExpressionStatement{Type: "ExpressionStatement", Expression: exp}
}

type AssignmentExpression struct {
	// assignment_expression ::= logical_or_expression [ assignment_operator assignment_expression ]
	// assignment_operator   ::= ASSIGN | PLUS_ASSIGN | MINUS_ASSIGN | STAR_ASSIGN | SLASH_ASSIGN
	Type     string     `json:"type"`
	Operator string     `json:"operator"`
	Left     Expression `json:"left"`
	Right    Expression `json:"right"`
}

func (ae *AssignmentExpression) expressionNode() {}
func (ae *AssignmentExpression) String() string {
	return fmt.Sprintf("(AssignmentExpression %s %v %v)", ae.Operator, ae.Left, ae.Right)
}

func NewAssignmentExpression(op string, left Expression, right Expression) *AssignmentExpression {
	return &AssignmentExpression{
		Type:     "AssignmentExpression",
		Operator: op,
		Left:     left,
		Right:    right,
	}
}

type LogicalExpression struct {
	// logical_or_expression   ::= logical_and_expression { OR logical_and_expression }
	// logical_and_expression   ::= equality_expression { AND equality_expression }
	Type     string     `json:"type"`
	Operator string     `json:"operator"`
	Left     Expression `json:"left"`
	Right    Expression `json:"right"`
}

func (le *LogicalExpression) expressionNode() {}
func (le *LogicalExpression) String() string {
	return fmt.Sprintf("(LogicalExpression %s %v %v)", le.Operator, le.Left, le.Right)
}

func NewLogicalExpression(op string, left, right Expression) *LogicalExpression {
	return &LogicalExpression{
		Type:     "LogicalExpression",
		Operator: op,
		Left:     left,
		Right:    right,
	}
}

type BinaryExpression struct {
	// equality_expression   		::= relational_expression { (EQ | NOT_EQ) relational_expression }
	// relational_expression 		::= additive_expression { (LT | LT_EQ | GT | GT_EQ) additive_expression }
	// additive_expression 			::= multiplicative_expression { (PLUS | MINUS) multiplicative_expression }
	// multiplicative_expression 	::= unary_expression { (STAR | SLASH | PERCENT) unary_expression }
	Type     string     `json:"type"`
	Operator string     `json:"operator"`
	Left     Expression `json:"left"`
	Right    Expression `json:"right"`
}

func (be *BinaryExpression) expressionNode() {}
func (be *BinaryExpression) String() string {
	return fmt.Sprintf("(BinaryExpression %s %v %v)", be.Operator, be.Left, be.Right)
}

func NewBinaryExpression(op string, left, right Expression) *BinaryExpression {
	return &BinaryExpression{
		Type:     "BinaryExpression",
		Operator: op,
		Left:     left,
		Right:    right,
	}
}

type UnaryExpression struct {
	// unary_expression	:= (MINUS | NOT) unary_expression | primary_expression
	Type     string     `json:"type"`
	Operator string     `json:"operator"`
	Right    Expression `json:"right"`
}

func (be *UnaryExpression) expressionNode() {}
func (be *UnaryExpression) String() string {
	return fmt.Sprintf("(UnaryExpression %s %v)", be.Operator, be.Right)
}

func NewUnaryExpression(op string, right Expression) *UnaryExpression {
	return &UnaryExpression{
		Type:     "UnaryExpression",
		Operator: op,
		Right:    right,
	}
}

type MemberExpression struct {
	// member_expression ::= (IDENT DOT IDENT | IDENT LBRACKET primary_expression RBRACKET)
	Type     string     `json:"type"`
	Computed bool       `json:"computed"`
	Object   Expression `json:"object"`
	Property Expression `json:"property"`
}

func (be *MemberExpression) expressionNode() {}
func (be *MemberExpression) String() string {
	return fmt.Sprintf("(MemberExpression %t %v %v)", be.Computed, be.Object, be.Property)
}

func NewMemberExpression(computed bool, object, property Expression) *MemberExpression {
	return &MemberExpression{
		Type:     "MemberExpression",
		Computed: computed,
		Object:   object,
		Property: property,
	}
}

type IntegerLiteral struct {
	// integer_literal ::= INT
	Type  string `json:"type"`
	Value int64  `json:"value"`
}

func (il *IntegerLiteral) expressionNode() {}

func (il *IntegerLiteral) String() string {
	return fmt.Sprintf("(IntegerLiteral %d)", il.Value)
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
func (fl *FloatLiteral) String() string {
	return fmt.Sprintf("(FloatLiteral %v)", fl.Value)
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
func (sl *StringLiteral) String() string {
	return fmt.Sprintf("(StringLiteral %s)", sl.Value)
}

func NewStringLiteral(value string) *StringLiteral {
	return &StringLiteral{Type: "StringLiteral", Value: value}
}

type BoolLiteral struct {
	// bool_literal ::= (TRUE | FALSE)
	Type  string `json:"type"`
	Value bool   `json:"value"`
}

func (bl *BoolLiteral) expressionNode() {}
func (bl *BoolLiteral) String() string {
	return fmt.Sprintf("(BoolLiteral %t)", bl.Value)
}

func NewBoolLiteral(value bool) *BoolLiteral {
	return &BoolLiteral{Type: "BoolLiteral", Value: value}
}

type NullLiteral struct {
	// null_literal ::= NULL
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func (nl *NullLiteral) expressionNode() {}
func (nl *NullLiteral) String() string {
	return fmt.Sprintf("(NullLiteral %v)", nl.Value)
}

func NewNullLiteral() *NullLiteral {
	return &NullLiteral{Type: "NullLiteral", Value: nil}
}

type Identifier struct {
	// identifier ::= IDENT
	Type string `json:"type"`
	Name string `json:"name"`
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string {
	return fmt.Sprintf("(Identifier %s)", i.Name)
}

func NewIdentifier(name string) *Identifier {
	return &Identifier{Type: "Identifier", Name: name}
}
