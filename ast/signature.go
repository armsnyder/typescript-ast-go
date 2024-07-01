package ast

// Signature is a [Node] that represents a signature. A signature defines a
// property.
type Signature interface {
	Expr
	signature()
}

// PropertySignature is an expression that defines an object property.
type PropertySignature struct {
	Name            *Identifier
	QuestionToken   bool
	Type            Type
	LeadingComment  string
	TrailingComment string
}

func (n *PropertySignature) String() string {
	return n.LeadingComment + " / " + n.TrailingComment
}

func (*PropertySignature) node()      {}
func (*PropertySignature) expr()      {}
func (*PropertySignature) signature() {}

// IndexSignature is an expression that defines an object index signature.
type IndexSignature struct {
	Parameters     []*Parameter
	Type           Type
	LeadingComment string
}

func (n *IndexSignature) String() string {
	return n.LeadingComment
}

func (*IndexSignature) node()      {}
func (*IndexSignature) expr()      {}
func (*IndexSignature) signature() {}
