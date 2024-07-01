package ast

// Type is a [Node] that represents a type expression. A type expression is a
// specific [Expr] that represents a type.
type Type interface {
	Expr
	typ()
}

// LiteralType is a literal type expression.
type LiteralType struct {
	Literal Expr
}

func (*LiteralType) node() {}
func (*LiteralType) expr() {}
func (*LiteralType) typ()  {}

// TypeLiteral is a type literal expression.
type TypeLiteral struct {
	Members []Signature
}

func (*TypeLiteral) node() {}
func (*TypeLiteral) expr() {}
func (*TypeLiteral) typ()  {}

// ArrayType is an array type expression.
type ArrayType struct {
	ElementType Expr
}

func (*ArrayType) node() {}
func (*ArrayType) expr() {}
func (*ArrayType) typ()  {}

// TypeReference is a type reference expression.
type TypeReference struct {
	TypeName Expr
}

func (*TypeReference) node() {}
func (*TypeReference) expr() {}
func (*TypeReference) typ()  {}

// UnionType is a union type expression.
type UnionType struct {
	Types []Type
}

func (*UnionType) node() {}
func (*UnionType) expr() {}
func (*UnionType) typ()  {}

// TupleType is a tuple type expression.
type TupleType struct {
	Elements []Type
}

func (*TupleType) node() {}
func (*TupleType) expr() {}
func (*TupleType) typ()  {}

// ParenthesizedType is an expression that wraps another expression in
// parentheses.
type ParenthesizedType struct {
	Type Type
}

func (*ParenthesizedType) node() {}
func (*ParenthesizedType) expr() {}
func (*ParenthesizedType) typ()  {}
