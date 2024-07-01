package ast_test

import (
	"fmt"

	"github.com/armsnyder/typescript-ast-go/ast"
	"github.com/armsnyder/typescript-ast-go/parser"
)

func Example() {
	sourceFile := parser.Parse([]byte(`
		export interface ProgressParams<T> {
		  token: ProgressToken;
		}`))

	ast.Inspect(sourceFile, func(node ast.Node) bool {
		if node != nil {
			fmt.Printf("Visited %T\n", node)
		}
		return true
	})

	// Output:
	// Visited *ast.SourceFile
	// Visited *ast.InterfaceDeclaration
	// Visited *ast.Identifier
	// Visited *ast.PropertySignature
	// Visited *ast.Identifier
	// Visited *ast.TypeReference
	// Visited *ast.Identifier
}
