package ast

import "fmt"

// Visitor is an interface for visiting nodes in the AST.
//
// The Visit method is called for each node in the AST. If the Visit method
// returns a non-nil Visitor, the children of the node are visited, followed by
// a call to w.Visit(nil).
type Visitor interface {
	Visit(node Node) (w Visitor)
}

// Walk traverses the AST rooted at node in depth-first order.
//
// It starts by calling v.Visit(node); node must not be nil. If the visitor
// returned by v.Visit(node) is not nil, Walk is called recursively with the
// visitor and each of the non-nil children of node, followed by a call to
// w.Visit(nil).
func Walk(v Visitor, node Node) { //nolint:revive // cyclomatic
	w := v.Visit(node)
	if w == nil {
		return
	}

	switch n := node.(type) {
	// Expressions.
	case *NumericLiteral, *StringLiteral, *Identifier:
	case *QualifiedName:
		Walk(w, n.Left)
		Walk(w, n.Right)
	case *ArrayLiteralExpression:
		for _, elem := range n.Elements {
			Walk(w, elem)
		}
	case *EnumMember:
		Walk(w, n.Name)
		Walk(w, n.Initializer)
	case *TypeParameter:
		Walk(w, n.Name)
	case *HeritageClause:
		for _, typ := range n.Types {
			Walk(w, typ)
		}
	case *ExpressionWithTypeArguments:
		Walk(w, n.Expression)
	case *PropertySignature:
		Walk(w, n.Name)
		Walk(w, n.Type)
	case *IndexSignature:
		for _, param := range n.Parameters {
			Walk(w, param)
		}
		Walk(w, n.Type)
	case *Parameter:
		Walk(w, n.Name)
		Walk(w, n.Type)
	case *VariableDeclarationList:
		for _, decl := range n.Declarations {
			Walk(w, decl)
		}
	case *VariableDeclaration:
		Walk(w, n.Name)
		Walk(w, n.Type)
		Walk(w, n.Initializer)
	case *PrefixUnaryExpression:
		Walk(w, n.Operand)

	// Types.
	case *LiteralType:
		Walk(w, n.Literal)
	case *TypeLiteral:
		for _, member := range n.Members {
			Walk(w, member)
		}
	case *ArrayType:
		Walk(w, n.ElementType)
	case *TypeReference:
		Walk(w, n.TypeName)
	case *UnionType:
		for _, typ := range n.Types {
			Walk(w, typ)
		}
	case *TupleType:
		for _, elem := range n.Elements {
			Walk(w, elem)
		}
	case *ParenthesizedType:
		Walk(w, n.Type)

	// Statements.
	case *SourceFile:
		for _, stmt := range n.Statements {
			Walk(w, stmt)
		}
	case *ModuleBlock:
		for _, stmt := range n.Statements {
			Walk(w, stmt)
		}
	case *VariableStatement:
		Walk(w, n.DeclarationList)
	case *TypeAliasDeclaration:
		Walk(w, n.Name)
		Walk(w, n.Type)
	case *EnumDeclaration:
		Walk(w, n.Name)
		for _, member := range n.Members {
			Walk(w, member)
		}
	case *InterfaceDeclaration:
		Walk(w, n.Name)
		for _, member := range n.Members {
			Walk(w, member)
		}
	case *ModuleDeclaration:
		Walk(w, n.Name)
		Walk(w, n.Body)

	default:
		panic(fmt.Sprintf("unknown node type %T", n))
	}

	w.Visit(nil)
}

type inspector func(Node) bool

func (f inspector) Visit(node Node) Visitor {
	if f(node) {
		return f
	}
	return nil
}

// Inspect traverses the AST rooted at node in depth-first order.
//
// It starts by calling f(node); node must not be nil. If f returns true,
// Inspect invokes f recursively for each of the non-nil children of node,
// followed by a call to f(nil).
func Inspect(node Node, f func(Node) bool) {
	Walk(inspector(f), node)
}
