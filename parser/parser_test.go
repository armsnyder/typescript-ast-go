package parser_test

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/armsnyder/typescript-ast-go/ast"
	"github.com/armsnyder/typescript-ast-go/parser"
	"github.com/armsnyder/typescript-ast-go/token"
)

func TestParser_Testdata(t *testing.T) {
	tests := []struct {
		testdata string
		want     *ast.SourceFile
	}{
		{
			testdata: "additional_properties",
			want: &ast.SourceFile{
				Statements: []ast.Stmt{
					&ast.InterfaceDeclaration{
						Name: &ast.Identifier{Text: "FormattingOptions"},
						Members: []ast.Signature{
							&ast.PropertySignature{
								Name: &ast.Identifier{Text: "tabSize"},
								Type: &ast.TypeReference{
									TypeName: &ast.Identifier{Text: "uinteger"},
								},
								LeadingComment: "Size of a tab in spaces.",
							},
							&ast.IndexSignature{
								Parameters: []*ast.Parameter{{
									Name: &ast.Identifier{Text: "key"},
									Type: &ast.TypeReference{
										TypeName: &ast.Identifier{Text: "string"},
									},
								}},
								Type: &ast.UnionType{
									Types: []ast.Type{
										&ast.TypeReference{
											TypeName: &ast.Identifier{Text: "boolean"},
										},
										&ast.TypeReference{
											TypeName: &ast.Identifier{Text: "integer"},
										},
										&ast.TypeReference{
											TypeName: &ast.Identifier{Text: "string"},
										},
									},
								},
								LeadingComment: "Signature for further properties.",
							},
						},
					},
				},
			},
		},
		{
			testdata: "array",
			want: &ast.SourceFile{
				Statements: []ast.Stmt{
					&ast.TypeAliasDeclaration{
						Name: &ast.Identifier{Text: "LSPArray"},
						Type: &ast.ArrayType{
							ElementType: &ast.TypeReference{
								TypeName: &ast.Identifier{Text: "LSPAny"},
							},
						},
						LeadingComment: "LSP arrays.\n\n@since 3.17.0",
					},
				},
			},
		},
		{
			testdata: "basic_type",
			want: &ast.SourceFile{
				Statements: []ast.Stmt{
					&ast.TypeAliasDeclaration{
						Name: &ast.Identifier{Text: "integer"},
						Type: &ast.TypeReference{
							TypeName: &ast.Identifier{Text: "number"},
						},
						LeadingComment: "Defines an integer number in the range of -2^31 to 2^31 - 1.",
					},
				},
			},
		},
		{
			testdata: "enum",
			want: &ast.SourceFile{
				Statements: []ast.Stmt{
					&ast.EnumDeclaration{
						Name: &ast.Identifier{Text: "SemanticTokenTypes"},
						Members: []*ast.EnumMember{
							{
								Name:        &ast.Identifier{Text: "namespace"},
								Initializer: &ast.StringLiteral{Text: "namespace"},
							},
							{
								Name:           &ast.Identifier{Text: "type"},
								Initializer:    &ast.StringLiteral{Text: "type"},
								LeadingComment: "Represents a generic type. Acts as a fallback for types which\ncan't be mapped to a specific type like class or enum.",
							},
							{
								Name:        &ast.Identifier{Text: "class"},
								Initializer: &ast.StringLiteral{Text: "class"},
							},
							{
								Name:        &ast.Identifier{Text: "enum"},
								Initializer: &ast.StringLiteral{Text: "enum"},
							},
							{
								Name:        &ast.Identifier{Text: "interface"},
								Initializer: &ast.StringLiteral{Text: "interface"},
							},
							{
								Name:        &ast.Identifier{Text: "string"},
								Initializer: &ast.StringLiteral{Text: "string"},
							},
						},
					},
				},
			},
		},
		{
			testdata: "generic",
			want: &ast.SourceFile{
				Statements: []ast.Stmt{
					&ast.InterfaceDeclaration{
						Name: &ast.Identifier{Text: "ProgressParams"},
						TypeParameters: []*ast.TypeParameter{
							{Name: &ast.Identifier{Text: "T"}},
						},
						Members: []ast.Signature{
							&ast.PropertySignature{
								Name: &ast.Identifier{Text: "token"},
								Type: &ast.TypeReference{
									TypeName: &ast.Identifier{Text: "ProgressToken"},
								},
								LeadingComment: "The progress token provided by the client or server.",
							},
							&ast.PropertySignature{
								Name: &ast.Identifier{Text: "value"},
								Type: &ast.TypeReference{
									TypeName: &ast.Identifier{Text: "T"},
								},
								LeadingComment: "The progress data.",
							},
						},
					},
				},
			},
		},
		{
			testdata: "inline_type",
			want: &ast.SourceFile{
				Statements: []ast.Stmt{
					&ast.InterfaceDeclaration{
						Name: &ast.Identifier{Text: "HoverParams"},
						Members: []ast.Signature{
							&ast.PropertySignature{
								Name: &ast.Identifier{Text: "textDocument"},
								Type: &ast.TypeReference{
									TypeName: &ast.Identifier{Text: "string"},
								},
								TrailingComment: "The text document's URI in string form",
							},
							&ast.PropertySignature{
								Name: &ast.Identifier{Text: "position"},
								Type: &ast.TypeLiteral{
									Members: []ast.Signature{
										&ast.PropertySignature{
											Name: &ast.Identifier{Text: "line"},
											Type: &ast.TypeReference{
												TypeName: &ast.Identifier{Text: "uinteger"},
											},
										},
										&ast.PropertySignature{
											Name: &ast.Identifier{Text: "character"},
											Type: &ast.TypeReference{
												TypeName: &ast.Identifier{Text: "uinteger"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			testdata: "interface",
			want: &ast.SourceFile{
				Statements: []ast.Stmt{
					&ast.InterfaceDeclaration{
						Name: &ast.Identifier{Text: "ResponseMessage"},
						HeritageClauses: []*ast.HeritageClause{
							{
								Types: []*ast.ExpressionWithTypeArguments{
									{Expression: &ast.Identifier{Text: "Message"}},
								},
							},
						},
						Members: []ast.Signature{
							&ast.PropertySignature{
								Name: &ast.Identifier{Text: "id"},
								Type: &ast.UnionType{
									Types: []ast.Type{
										&ast.TypeReference{TypeName: &ast.Identifier{Text: "integer"}},
										&ast.TypeReference{TypeName: &ast.Identifier{Text: "string"}},
										&ast.TypeReference{TypeName: &ast.Identifier{Text: "null"}},
									},
								},
								LeadingComment: "The request id.",
							},
							&ast.PropertySignature{
								Name:          &ast.Identifier{Text: "result"},
								QuestionToken: true,
								Type: &ast.UnionType{
									Types: []ast.Type{
										&ast.TypeReference{TypeName: &ast.Identifier{Text: "string"}},
										&ast.TypeReference{TypeName: &ast.Identifier{Text: "number"}},
										&ast.TypeReference{TypeName: &ast.Identifier{Text: "boolean"}},
										&ast.TypeReference{TypeName: &ast.Identifier{Text: "array"}},
										&ast.TypeReference{TypeName: &ast.Identifier{Text: "object"}},
										&ast.TypeReference{TypeName: &ast.Identifier{Text: "null"}},
									},
								},
								LeadingComment: "The result of a request. This member is REQUIRED on success.\nThis member MUST NOT exist if there was an error invoking the method.",
							},
							&ast.PropertySignature{
								Name:          &ast.Identifier{Text: "error"},
								QuestionToken: true,
								Type: &ast.TypeReference{
									TypeName: &ast.Identifier{Text: "ResponseError"},
								},
								LeadingComment: "The error object in case a request fails.",
							},
						},
					},
				},
			},
		},
		{
			testdata: "interface_complex_syntax",
			want: &ast.SourceFile{
				Statements: []ast.Stmt{
					&ast.InterfaceDeclaration{
						Name: &ast.Identifier{Text: "WorkspaceEdit"},
						Members: []ast.Signature{
							&ast.PropertySignature{
								Name:          &ast.Identifier{Text: "changes"},
								QuestionToken: true,
								Type: &ast.TypeLiteral{
									Members: []ast.Signature{
										&ast.IndexSignature{
											Parameters: []*ast.Parameter{{
												Name: &ast.Identifier{Text: "uri"},
												Type: &ast.TypeReference{
													TypeName: &ast.Identifier{Text: "DocumentUri"},
												},
											}},
											Type: &ast.ArrayType{
												ElementType: &ast.TypeReference{
													TypeName: &ast.Identifier{Text: "TextEdit"},
												},
											},
										},
									},
								},
								LeadingComment: "Holds changes to existing resources.",
							},
							&ast.PropertySignature{
								Name:          &ast.Identifier{Text: "documentChanges"},
								QuestionToken: true,
								Type: &ast.ParenthesizedType{
									Type: &ast.UnionType{
										Types: []ast.Type{
											&ast.ArrayType{
												ElementType: &ast.TypeReference{
													TypeName: &ast.Identifier{Text: "TextDocumentEdit"},
												},
											},
											&ast.ArrayType{
												ElementType: &ast.ParenthesizedType{
													Type: &ast.UnionType{
														Types: []ast.Type{
															&ast.TypeReference{
																TypeName: &ast.Identifier{Text: "TextDocumentEdit"},
															},
															&ast.TypeReference{
																TypeName: &ast.Identifier{Text: "CreateFile"},
															},
															&ast.TypeReference{
																TypeName: &ast.Identifier{Text: "RenameFile"},
															},
															&ast.TypeReference{
																TypeName: &ast.Identifier{Text: "DeleteFile"},
															},
														},
													},
												},
											},
										},
									},
								},
								LeadingComment: strings.TrimSpace(`
Depending on the client capability
` + "`workspace.workspaceEdit.resourceOperations`" + ` document changes are either
an array of ` + "`TextDocumentEdit`" + `s to express changes to n different text
documents where each text document edit addresses a specific version of
a text document. Or it can contain above ` + "`TextDocumentEdit`" + `s mixed with
create, rename and delete file / folder operations.

Whether a client supports versioned document edits is expressed via
` + "`workspace.workspaceEdit.documentChanges`" + ` client capability.

If a client neither supports ` + "`documentChanges`" + ` nor
` + "`workspace.workspaceEdit.resourceOperations`" + ` then only plain ` + "`TextEdit`" + `s
using the ` + "`changes`" + ` property are supported.
`),
							},
							&ast.PropertySignature{
								Name:          &ast.Identifier{Text: "changeAnnotations"},
								QuestionToken: true,
								Type: &ast.TypeLiteral{
									Members: []ast.Signature{
										&ast.IndexSignature{
											Parameters: []*ast.Parameter{{
												Name: &ast.Identifier{Text: "id"},
												Type: &ast.TypeReference{
													TypeName: &ast.Identifier{Text: "string"},
												},
											}},
											Type: &ast.TypeReference{
												TypeName: &ast.Identifier{Text: "ChangeAnnotation"},
											},
										},
									},
								},
								LeadingComment: strings.TrimSpace(`
A map of change annotations that can be referenced in
` + "`AnnotatedTextEdit`" + `s or create, rename and delete file / folder
operations.

Whether clients honor this property depends on the client capability
` + "`workspace.changeAnnotationSupport`" + `.

@since 3.16.0
`),
							},
						},
					},
				},
			},
		},
		{
			testdata: "interface_with_union_array",
			want: &ast.SourceFile{
				Statements: []ast.Stmt{
					&ast.InterfaceDeclaration{
						Name: &ast.Identifier{Text: "TextDocumentEdit"},
						Members: []ast.Signature{
							&ast.PropertySignature{
								Name: &ast.Identifier{Text: "edits"},
								Type: &ast.ArrayType{
									ElementType: &ast.ParenthesizedType{
										Type: &ast.UnionType{
											Types: []ast.Type{
												&ast.TypeReference{
													TypeName: &ast.Identifier{Text: "TextEdit"},
												},
												&ast.TypeReference{
													TypeName: &ast.Identifier{Text: "AnnotatedTextEdit"},
												},
											},
										},
									},
								},
								LeadingComment: strings.TrimSpace(`
The edits to be applied.

@since 3.16.0 - support for AnnotatedTextEdit. This is guarded by the
client capability ` + "`workspace.workspaceEdit.changeAnnotationSupport`"),
							},
						},
					},
				},
			},
		},
		{
			testdata: "interface_with_union_struct_field",
			want: &ast.SourceFile{
				Statements: []ast.Stmt{
					&ast.InterfaceDeclaration{
						Name: &ast.Identifier{Text: "NotebookDocumentSyncOptions"},
						Members: []ast.Signature{
							&ast.PropertySignature{
								Name: &ast.Identifier{Text: "notebookSelector"},
								Type: &ast.ArrayType{
									ElementType: &ast.ParenthesizedType{
										Type: &ast.UnionType{
											Types: []ast.Type{
												&ast.TypeLiteral{
													Members: []ast.Signature{
														&ast.PropertySignature{
															Name: &ast.Identifier{Text: "notebook"},
															Type: &ast.UnionType{
																Types: []ast.Type{
																	&ast.TypeReference{
																		TypeName: &ast.Identifier{Text: "string"},
																	},
																	&ast.TypeReference{
																		TypeName: &ast.Identifier{Text: "NotebookDocumentFilter"},
																	},
																},
															},
															LeadingComment: "The notebook to be synced. If a string\nvalue is provided it matches against the\nnotebook type. '*' matches every notebook.",
														},
														&ast.PropertySignature{
															Name:          &ast.Identifier{Text: "cells"},
															QuestionToken: true,
															Type: &ast.ArrayType{
																ElementType: &ast.TypeLiteral{
																	Members: []ast.Signature{
																		&ast.PropertySignature{
																			Name: &ast.Identifier{Text: "language"},
																			Type: &ast.TypeReference{
																				TypeName: &ast.Identifier{Text: "string"},
																			},
																		},
																	},
																},
															},
															LeadingComment: "The cells of the matching notebook to be synced.",
														},
													},
												},
												&ast.TypeLiteral{
													Members: []ast.Signature{
														&ast.PropertySignature{
															Name:          &ast.Identifier{Text: "notebook"},
															QuestionToken: true,
															Type: &ast.UnionType{
																Types: []ast.Type{
																	&ast.TypeReference{
																		TypeName: &ast.Identifier{Text: "string"},
																	},
																	&ast.TypeReference{
																		TypeName: &ast.Identifier{Text: "NotebookDocumentFilter"},
																	},
																},
															},
															LeadingComment: "The notebook to be synced. If a string\nvalue is provided it matches against the\nnotebook type. '*' matches every notebook.",
														},
														&ast.PropertySignature{
															Name: &ast.Identifier{Text: "cells"},
															Type: &ast.ArrayType{
																ElementType: &ast.TypeLiteral{
																	Members: []ast.Signature{
																		&ast.PropertySignature{
																			Name: &ast.Identifier{Text: "language"},
																			Type: &ast.TypeReference{
																				TypeName: &ast.Identifier{Text: "string"},
																			},
																		},
																	},
																},
															},
															LeadingComment: "The cells of the matching notebook to be synced.",
														},
													},
												},
											},
										},
									},
								},
								LeadingComment: "The notebooks to be synced",
							},
							&ast.PropertySignature{
								Name:          &ast.Identifier{Text: "save"},
								QuestionToken: true,
								Type: &ast.TypeReference{
									TypeName: &ast.Identifier{Text: "boolean"},
								},
								LeadingComment: "Whether save notification should be forwarded to\nthe server. Will only be honored if mode === `notebook`.",
							},
						},
						LeadingComment: strings.TrimSpace(`
Options specific to a notebook plus its cells
to be synced to the server.

If a selector provides a notebook document
filter but no cell selector all cells of a
matching notebook document will be synced.

If a selector provides no notebook document
filter but only a cell selector all notebook
documents that contain at least one matching
cell will be synced.

@since 3.17.0`),
					},
				},
			},
		},
		{
			testdata: "multiple_extends",
			want: &ast.SourceFile{
				Statements: []ast.Stmt{
					&ast.InterfaceDeclaration{
						Name: &ast.Identifier{Text: "NotebookDocumentSyncRegistrationOptions"},
						HeritageClauses: []*ast.HeritageClause{
							{
								Types: []*ast.ExpressionWithTypeArguments{
									{Expression: &ast.Identifier{Text: "NotebookDocumentSyncOptions"}},
								},
							},
							{
								Types: []*ast.ExpressionWithTypeArguments{
									{Expression: &ast.Identifier{Text: "StaticRegistrationOptions"}},
								},
							},
						},
						LeadingComment: "Registration options specific to a notebook.\n\n@since 3.17.0",
					},
				},
			},
		},
		{
			testdata: "namespace",
			want: &ast.SourceFile{
				Statements: []ast.Stmt{
					&ast.ModuleDeclaration{
						Name: &ast.Identifier{Text: "ErrorCodes"},
						Body: &ast.ModuleBlock{
							Statements: []ast.Stmt{
								&ast.VariableStatement{
									DeclarationList: &ast.VariableDeclarationList{
										Declarations: []*ast.VariableDeclaration{{
											Name: &ast.Identifier{Text: "ParseError"},
											Type: &ast.TypeReference{
												TypeName: &ast.Identifier{Text: "integer"},
											},
											Initializer: &ast.PrefixUnaryExpression{
												Operator: token.Minus,
												Operand:  &ast.NumericLiteral{Text: "32700"},
											},
										}},
									},
									LeadingComment: "Defined by JSON-RPC",
								},
								&ast.VariableStatement{
									DeclarationList: &ast.VariableDeclarationList{
										Declarations: []*ast.VariableDeclaration{{
											Name: &ast.Identifier{Text: "InvalidRequest"},
											Type: &ast.TypeReference{
												TypeName: &ast.Identifier{Text: "integer"},
											},
											Initializer: &ast.PrefixUnaryExpression{
												Operator: token.Minus,
												Operand:  &ast.NumericLiteral{Text: "32600"},
											},
										}},
									},
								},
								&ast.VariableStatement{
									DeclarationList: &ast.VariableDeclarationList{
										Declarations: []*ast.VariableDeclaration{{
											Name: &ast.Identifier{Text: "jsonrpcReservedErrorRangeStart"},
											Type: &ast.TypeReference{
												TypeName: &ast.Identifier{Text: "integer"},
											},
											Initializer: &ast.PrefixUnaryExpression{
												Operator: token.Minus,
												Operand:  &ast.NumericLiteral{Text: "32099"},
											},
										}},
									},
									LeadingComment: strings.TrimSpace(`
This is the start range of JSON-RPC reserved error codes.
It doesn't denote a real error code. No LSP error codes should
be defined between the start and end range. For backwards
compatibility the ` + "`ServerNotInitialized`" + ` and the ` + "`UnknownErrorCode`" + `
are left in the range.

@since 3.16.0
`),
								},
								&ast.VariableStatement{
									DeclarationList: &ast.VariableDeclarationList{
										Declarations: []*ast.VariableDeclaration{{
											Name: &ast.Identifier{Text: "serverErrorStart"},
											Type: &ast.TypeReference{
												TypeName: &ast.Identifier{Text: "integer"},
											},
											Initializer: &ast.TypeReference{
												TypeName: &ast.Identifier{Text: "jsonrpcReservedErrorRangeStart"},
											},
										}},
									},
									LeadingComment: "@deprecated use jsonrpcReservedErrorRangeStart",
								},
							},
						},
					},
				},
			},
		},
		{
			testdata: "namespace_with_type",
			want: &ast.SourceFile{
				Statements: []ast.Stmt{
					&ast.ModuleDeclaration{
						Name: &ast.Identifier{Text: "DiagnosticSeverity"},
						Body: &ast.ModuleBlock{
							Statements: []ast.Stmt{
								&ast.VariableStatement{
									DeclarationList: &ast.VariableDeclarationList{
										Declarations: []*ast.VariableDeclaration{{
											Name: &ast.Identifier{Text: "Error"},
											Type: &ast.LiteralType{
												Literal: &ast.NumericLiteral{Text: "1"},
											},
											Initializer: &ast.NumericLiteral{Text: "1"},
										}},
									},
									LeadingComment: "Reports an error.",
								},
								&ast.VariableStatement{
									DeclarationList: &ast.VariableDeclarationList{
										Declarations: []*ast.VariableDeclaration{{
											Name: &ast.Identifier{Text: "Warning"},
											Type: &ast.LiteralType{
												Literal: &ast.NumericLiteral{Text: "2"},
											},
											Initializer: &ast.NumericLiteral{Text: "2"},
										}},
									},
									LeadingComment: "Reports a warning.",
								},
							},
						},
					},
					&ast.TypeAliasDeclaration{
						Name: &ast.Identifier{Text: "DiagnosticSeverity"},
						Type: &ast.UnionType{
							Types: []ast.Type{
								&ast.LiteralType{
									Literal: &ast.NumericLiteral{Text: "1"},
								},
								&ast.LiteralType{
									Literal: &ast.NumericLiteral{Text: "2"},
								},
							},
						},
					},
				},
			},
		},
		{
			testdata: "object",
			want: &ast.SourceFile{
				Statements: []ast.Stmt{
					&ast.TypeAliasDeclaration{
						Name: &ast.Identifier{Text: "LSPObject"},
						Type: &ast.TypeLiteral{
							Members: []ast.Signature{
								&ast.IndexSignature{
									Parameters: []*ast.Parameter{{
										Name: &ast.Identifier{Text: "key"},
										Type: &ast.TypeReference{
											TypeName: &ast.Identifier{Text: "string"},
										},
									}},
									Type: &ast.TypeReference{
										TypeName: &ast.Identifier{Text: "LSPAny"},
									},
								},
							},
						},
						LeadingComment: "LSP object definition.\n\n@since 3.17.0",
					},
				},
			},
		},
		{
			testdata: "qualified",
			want: &ast.SourceFile{
				Statements: []ast.Stmt{
					&ast.InterfaceDeclaration{
						Name: &ast.Identifier{Text: "FullDocumentDiagnosticReport"},
						Members: []ast.Signature{
							&ast.PropertySignature{
								Name: &ast.Identifier{Text: "kind"},
								Type: &ast.TypeReference{
									TypeName: &ast.QualifiedName{
										Left:  &ast.Identifier{Text: "DocumentDiagnosticReportKind"},
										Right: &ast.Identifier{Text: "Full"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			testdata: "readonly",
			want: &ast.SourceFile{
				Statements: []ast.Stmt{
					&ast.InterfaceDeclaration{
						Name: &ast.Identifier{Text: "SemanticTokensDelta"},
						Members: []ast.Signature{
							&ast.PropertySignature{
								Name:          &ast.Identifier{Text: "resultId"},
								QuestionToken: true,
								Type: &ast.TypeReference{
									TypeName: &ast.Identifier{Text: "string"},
								},
							},
							&ast.PropertySignature{
								Name: &ast.Identifier{Text: "edits"},
								Type: &ast.ArrayType{
									ElementType: &ast.TypeReference{
										TypeName: &ast.Identifier{Text: "SemanticTokensEdit"},
									},
								},
								LeadingComment: "The semantic token edits to transform a previous result into a new\nresult.",
							},
						},
					},
				},
			},
		},
		{
			testdata: "tuple",
			want: &ast.SourceFile{
				Statements: []ast.Stmt{
					&ast.InterfaceDeclaration{
						Name: &ast.Identifier{Text: "ParameterInformation"},
						Members: []ast.Signature{
							&ast.PropertySignature{
								Name: &ast.Identifier{Text: "label"},
								Type: &ast.UnionType{
									Types: []ast.Type{
										&ast.TypeReference{
											TypeName: &ast.Identifier{Text: "string"},
										},
										&ast.TupleType{
											Elements: []ast.Type{
												&ast.TypeReference{
													TypeName: &ast.Identifier{Text: "uinteger"},
												},
												&ast.TypeReference{
													TypeName: &ast.Identifier{Text: "uinteger"},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			testdata: "union_struct",
			want: &ast.SourceFile{
				Statements: []ast.Stmt{
					&ast.TypeAliasDeclaration{
						Name: &ast.Identifier{Text: "NotebookDocumentFilter"},
						Type: &ast.UnionType{
							Types: []ast.Type{
								&ast.TypeLiteral{
									Members: []ast.Signature{
										&ast.PropertySignature{
											Name: &ast.Identifier{Text: "notebookType"},
											Type: &ast.TypeReference{
												TypeName: &ast.Identifier{Text: "string"},
											},
											LeadingComment: "The type of the enclosing notebook.",
										},
										&ast.PropertySignature{
											Name:          &ast.Identifier{Text: "scheme"},
											QuestionToken: true,
											Type: &ast.TypeReference{
												TypeName: &ast.Identifier{Text: "string"},
											},
											LeadingComment: "A Uri [scheme](#Uri.scheme), like `file` or `untitled`.",
										},
										&ast.PropertySignature{
											Name:          &ast.Identifier{Text: "pattern"},
											QuestionToken: true,
											Type: &ast.TypeReference{
												TypeName: &ast.Identifier{Text: "string"},
											},
											LeadingComment: "A glob pattern.",
										},
									},
								},
								&ast.TypeLiteral{
									Members: []ast.Signature{
										&ast.PropertySignature{
											Name:          &ast.Identifier{Text: "notebookType"},
											QuestionToken: true,
											Type: &ast.TypeReference{
												TypeName: &ast.Identifier{Text: "string"},
											},
											LeadingComment: "The type of the enclosing notebook.",
										},
										&ast.PropertySignature{
											Name: &ast.Identifier{Text: "scheme"},
											Type: &ast.TypeReference{
												TypeName: &ast.Identifier{Text: "string"},
											},
											LeadingComment: "A Uri [scheme](#Uri.scheme), like `file` or `untitled`.",
										},
										&ast.PropertySignature{
											Name:          &ast.Identifier{Text: "pattern"},
											QuestionToken: true,
											Type: &ast.TypeReference{
												TypeName: &ast.Identifier{Text: "string"},
											},
											LeadingComment: "A glob pattern.",
										},
									},
								},
								&ast.TypeLiteral{
									Members: []ast.Signature{
										&ast.PropertySignature{
											Name:          &ast.Identifier{Text: "notebookType"},
											QuestionToken: true,
											Type: &ast.TypeReference{
												TypeName: &ast.Identifier{Text: "string"},
											},
											LeadingComment: "The type of the enclosing notebook.",
										},
										&ast.PropertySignature{
											Name:          &ast.Identifier{Text: "scheme"},
											QuestionToken: true,
											Type: &ast.TypeReference{
												TypeName: &ast.Identifier{Text: "string"},
											},
											LeadingComment: "A Uri [scheme](#Uri.scheme), like `file` or `untitled`.",
										},
										&ast.PropertySignature{
											Name: &ast.Identifier{Text: "pattern"},
											Type: &ast.TypeReference{
												TypeName: &ast.Identifier{Text: "string"},
											},
											LeadingComment: "A glob pattern.",
										},
									},
								},
							},
						},
						LeadingComment: "A notebook document filter denotes a notebook document by\ndifferent properties.\n\n@since 3.17.0",
					},
				},
			},
		},
		{
			testdata: "union_type",
			want: &ast.SourceFile{
				Statements: []ast.Stmt{
					&ast.TypeAliasDeclaration{
						Name: &ast.Identifier{Text: "LSPAny"},
						Type: &ast.UnionType{
							Types: []ast.Type{
								&ast.TypeReference{TypeName: &ast.Identifier{Text: "LSPObject"}},
								&ast.TypeReference{TypeName: &ast.Identifier{Text: "LSPArray"}},
								&ast.TypeReference{TypeName: &ast.Identifier{Text: "string"}},
								&ast.TypeReference{TypeName: &ast.Identifier{Text: "integer"}},
								&ast.TypeReference{TypeName: &ast.Identifier{Text: "uinteger"}},
								&ast.TypeReference{TypeName: &ast.Identifier{Text: "decimal"}},
								&ast.TypeReference{TypeName: &ast.Identifier{Text: "boolean"}},
								&ast.TypeReference{TypeName: &ast.Identifier{Text: "null"}},
							},
						},
						LeadingComment: "The LSP any type\n\n@since 3.17.0",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testdata, func(t *testing.T) {
			source, err := os.ReadFile(fmt.Sprintf("../internal/testdata/%s.ts.txt", tt.testdata))
			if err != nil {
				t.Fatal(err)
			}

			if got := parser.Parse(source); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot:\n%s\n\nwant:\n%s", printTreeStructure(got), printTreeStructure(tt.want))
			}
		})
	}
}

func TestParser_Inline(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want *ast.SourceFile
	}{
		{
			name: "array",
			src:  `export const EOL: string[] = ['\n', '\r\n', '\r'];`,
			want: &ast.SourceFile{
				Statements: []ast.Stmt{
					&ast.VariableStatement{
						DeclarationList: &ast.VariableDeclarationList{
							Declarations: []*ast.VariableDeclaration{{
								Name: &ast.Identifier{Text: "EOL"},
								Type: &ast.ArrayType{
									ElementType: &ast.TypeReference{
										TypeName: &ast.Identifier{Text: "string"},
									},
								},
								Initializer: &ast.ArrayLiteralExpression{
									Elements: []ast.Expr{
										&ast.StringLiteral{Text: `\n`},
										&ast.StringLiteral{Text: `\r\n`},
										&ast.StringLiteral{Text: `\r`},
									},
								},
							}},
						},
					},
				},
			},
		},
		{
			name: "identifier with numbers",
			src:  `export const UTF8: string = 'utf-8';`,
			want: &ast.SourceFile{
				Statements: []ast.Stmt{
					&ast.VariableStatement{
						DeclarationList: &ast.VariableDeclarationList{
							Declarations: []*ast.VariableDeclaration{{
								Name: &ast.Identifier{Text: "UTF8"},
								Type: &ast.TypeReference{
									TypeName: &ast.Identifier{Text: "string"},
								},
								Initializer: &ast.StringLiteral{Text: "utf-8"},
							}},
						},
					},
				},
			},
		},
		{
			name: "keyword property",
			src: `
interface FileEvent {
	type: FileChangeType;
}`,
			want: &ast.SourceFile{
				Statements: []ast.Stmt{
					&ast.InterfaceDeclaration{
						Name: &ast.Identifier{Text: "FileEvent"},
						Members: []ast.Signature{
							&ast.PropertySignature{
								Name: &ast.Identifier{Text: "type"},
								Type: &ast.TypeReference{
									TypeName: &ast.Identifier{Text: "FileChangeType"},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parser.Parse([]byte(tt.src)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot:\n%s\n\nwant:\n%s", printTreeStructure(got), printTreeStructure(tt.want))
			}
		})
	}
}

func printTreeStructure(node ast.Node) string {
	if node == nil {
		return ""
	}

	sb := &strings.Builder{}
	indent := 0

	ast.Inspect(node, func(node ast.Node) bool {
		if node == nil {
			indent--
			return false
		}

		for i := 0; i < indent; i++ {
			sb.WriteString("  ")
		}

		sb.WriteString(reflect.TypeOf(node).String())
		if stringer, ok := node.(fmt.Stringer); ok {
			sb.WriteString(" ")
			sb.WriteString(stringer.String())
		}
		sb.WriteString("\n")

		indent++
		return true
	})

	return sb.String()
}
