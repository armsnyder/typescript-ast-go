# TypeScript AST

This library provides a way to parse TypeScript source code into an abstract
syntax tree (AST) in Golang. The packages are laid out similar to the standard
[go](https://pkg.go.dev/go) library.

[Package Documentation](https://pkg.go.dev/github.com/armsnyder/typescript-ast-go)

The main two packages are:

- [parser](https://pkg.go.dev/github.com/armsnyder/typescript-ast-go/parser):
  Parse TypeScript source code into an AST.
- [ast](https://pkg.go.dev/github.com/armsnyder/typescript-ast-go/ast): The
  AST nodes and visitor for TypeScript source code.

This library was originally created in order to parse TypeScript type
definitions specifically for the Language Server Protocol Specification. As a
result, it is not feature complete and may not work for all TypeScript source
code.
