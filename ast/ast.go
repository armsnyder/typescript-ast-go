// Package ast defines the abstract syntax tree (AST) for the TypeScript
// programming language and provides functionality for traversing the AST.
package ast

// Node is a common interface that all nodes in the AST implement.
type Node interface {
	node()
}
