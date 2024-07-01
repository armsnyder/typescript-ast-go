package parser_test

import (
	"fmt"

	"github.com/armsnyder/typescript-ast-go/parser"
)

func Example() {
	src := []byte(`
		export interface ProgressParams<T> {
		  /**
		   * The progress token provided by the client or server.
		   */
		  token: ProgressToken;

		  /**
		   * The progress data.
		   */
		  value: T;
		}`)
	sourceFile := parser.Parse(src)
	fmt.Printf("Parsed %T", sourceFile.Statements[0])
	// Output: Parsed *ast.InterfaceDeclaration
}
