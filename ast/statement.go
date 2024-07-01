package ast

// Stmt is a [Node] that represents a statement. A statement performs an
// action.
type Stmt interface {
	Node
	stmt()
}

// SourceFile is a statement that represents a source file.
type SourceFile struct {
	Statements []Stmt
}

func (*SourceFile) node() {}
func (*SourceFile) stmt() {}

// ModuleBlock is a statement that represents a block of statements in a
// module.
type ModuleBlock struct {
	Statements []Stmt
}

func (*ModuleBlock) node() {}
func (*ModuleBlock) stmt() {}

// VariableStatement is a statement that declares a variable.
type VariableStatement struct {
	DeclarationList *VariableDeclarationList
	LeadingComment  string
}

func (n *VariableStatement) String() string {
	return n.LeadingComment
}

func (*VariableStatement) node() {}
func (*VariableStatement) stmt() {}

// TypeAliasDeclaration is a statement that introduces a new type alias.
type TypeAliasDeclaration struct {
	Name           *Identifier
	Type           Type
	LeadingComment string
}

func (n *TypeAliasDeclaration) String() string {
	return n.LeadingComment
}

func (*TypeAliasDeclaration) node() {}
func (*TypeAliasDeclaration) stmt() {}

// EnumDeclaration is a statement that introduces a new enum.
type EnumDeclaration struct {
	Name           *Identifier
	Members        []*EnumMember
	LeadingComment string
}

func (n *EnumDeclaration) String() string {
	return n.LeadingComment
}

func (*EnumDeclaration) node() {}
func (*EnumDeclaration) stmt() {}

// InterfaceDeclaration is a statement that introduces a new interface.
type InterfaceDeclaration struct {
	Name            *Identifier
	TypeParameters  []*TypeParameter
	HeritageClauses []*HeritageClause
	Members         []Signature
	LeadingComment  string
}

func (n *InterfaceDeclaration) String() string {
	return n.LeadingComment
}

func (*InterfaceDeclaration) node() {}
func (*InterfaceDeclaration) stmt() {}

type ModuleDeclaration struct {
	Name           *Identifier
	Body           *ModuleBlock
	LeadingComment string
}

func (*ModuleDeclaration) node() {}
func (*ModuleDeclaration) stmt() {}
