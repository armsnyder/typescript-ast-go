package ast

import "github.com/armsnyder/typescript-ast-go/token"

// Expr is a [Node] that represents an expression. An expression produces a
// value.
type Expr interface {
	Node
	expr()
}

// NumericLiteral is a numeric literal expression.
type NumericLiteral struct {
	Text string
}

func (n *NumericLiteral) String() string {
	return n.Text
}

func (*NumericLiteral) node() {}
func (*NumericLiteral) expr() {}

// StringLiteral is a string literal expression.
type StringLiteral struct {
	Text string
}

func (n *StringLiteral) String() string {
	return n.Text
}

func (*StringLiteral) node() {}
func (*StringLiteral) expr() {}

// ArrayLiteralExpression is an array literal expression.
type ArrayLiteralExpression struct {
	Elements []Expr
}

func (*ArrayLiteralExpression) node() {}
func (*ArrayLiteralExpression) expr() {}

// Identifier is an identifier literal expression.
type Identifier struct {
	Text string
}

func (n *Identifier) String() string {
	return n.Text
}

func (*Identifier) node() {}
func (*Identifier) expr() {}

// QualifiedName is a qualified name expression.
type QualifiedName struct {
	Left  *Identifier
	Right *Identifier
}

func (*QualifiedName) node() {}
func (*QualifiedName) expr() {}

// EnumMember is an enum member expression.
type EnumMember struct {
	Name           *Identifier
	Initializer    Expr
	LeadingComment string
}

func (n *EnumMember) String() string {
	return n.LeadingComment
}

func (*EnumMember) node() {}
func (*EnumMember) expr() {}

// TypeParameter is a type parameter expression.
type TypeParameter struct {
	Name *Identifier
}

func (*TypeParameter) node() {}
func (*TypeParameter) expr() {}

// HeritageClause is a heritage clause expression.
type HeritageClause struct {
	Types []*ExpressionWithTypeArguments
}

func (*HeritageClause) node() {}
func (*HeritageClause) expr() {}

// ExpressionWithTypeArguments is an expression with type arguments.
type ExpressionWithTypeArguments struct {
	Expression *Identifier
}

func (*ExpressionWithTypeArguments) node() {}
func (*ExpressionWithTypeArguments) expr() {}

// Parameter is a parameter expression.
type Parameter struct {
	Name *Identifier
	Type Type
}

func (*Parameter) node() {}
func (*Parameter) expr() {}

// VariableDeclarationList is an expression that declares a list of variables.
type VariableDeclarationList struct {
	Declarations []*VariableDeclaration
}

func (*VariableDeclarationList) node() {}
func (*VariableDeclarationList) expr() {}

// VariableDeclaration is an expression that declares a variable.
type VariableDeclaration struct {
	Name        *Identifier
	Type        Type
	Initializer Expr
}

func (*VariableDeclaration) node() {}
func (*VariableDeclaration) expr() {}

// PrefixUnaryExpression is an expression that applies a unary operator to an
// operand.
type PrefixUnaryExpression struct {
	Operator token.Kind
	Operand  Expr
}

func (*PrefixUnaryExpression) node() {}
func (*PrefixUnaryExpression) expr() {}
