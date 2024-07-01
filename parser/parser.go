// Package parser provides a [Parse] function for parsing TypeScript source
// files into an abstract syntax tree (AST).
package parser

import (
	"fmt"

	"github.com/armsnyder/typescript-ast-go/ast"
	"github.com/armsnyder/typescript-ast-go/token"
)

func Parse(source []byte) *ast.SourceFile {
	p := parser{lex: &lexer{Source: source}}
	return p.parseSourceFile()
}

type parser struct {
	lex             *lexer
	tok             token.Token
	lastComment     string
	lastLineComment string
}

func (p *parser) parseSourceFile() *ast.SourceFile {
	sourceFile := &ast.SourceFile{}

	for {
		p.advance()
		if p.tok.Kind == token.EOF {
			return sourceFile
		}
		sourceFile.Statements = append(sourceFile.Statements, p.parseStatement())
	}
}

func (p *parser) parseStatement() ast.Stmt {
	p.expect(token.Ident)
	for {
		switch p.tok.Text {
		case "export":
			p.advance()
		case "const":
			return p.parseVariableStatement()
		case "type":
			return p.parseTypeAliasDeclaration()
		case "enum":
			return p.parseEnumDeclaration()
		case "interface":
			return p.parseInterfaceDeclaration()
		case "namespace":
			return p.parseModuleDeclaration()
		default:
			panic(fmt.Sprintf("unexpected token %s", p.tok))
		}
	}
}

func (p *parser) parseVariableStatement() *ast.VariableStatement {
	p.eat(token.Ident)
	decl := &ast.VariableStatement{LeadingComment: p.consumeComment()}
	decl.DeclarationList = p.parseVariableDeclarationList()
	return decl
}

func (p *parser) parseVariableDeclarationList() *ast.VariableDeclarationList {
	decl := &ast.VariableDeclarationList{}
	for {
		decl.Declarations = append(decl.Declarations, p.parseVariableDeclaration())
		if p.tok.Kind != token.Comma {
			p.eat(token.Semicolon)
			return decl
		}
		p.advance()
	}
}

func (p *parser) parseVariableDeclaration() *ast.VariableDeclaration {
	decl := &ast.VariableDeclaration{}
	decl.Name = p.parseIdentifier()
	if p.tok.Kind == token.Colon {
		p.advance()
		decl.Type = p.parseType()
	}
	if p.tok.Kind == token.Assign {
		p.advance()
		decl.Initializer = p.parseInitializer()
	}
	return decl
}

func (p *parser) parseModuleDeclaration() *ast.ModuleDeclaration {
	p.eat(token.Ident)
	decl := &ast.ModuleDeclaration{LeadingComment: p.consumeComment()}
	decl.Name = p.parseIdentifier()
	p.eat(token.LBrace)
	decl.Body = &ast.ModuleBlock{}
	for {
		if p.tok.Kind == token.RBrace {
			p.eat(token.RBrace)
			return decl
		}
		decl.Body.Statements = append(decl.Body.Statements, p.parseStatement())
	}
}

func (p *parser) parseTypeAliasDeclaration() *ast.TypeAliasDeclaration {
	p.eat(token.Ident)
	decl := &ast.TypeAliasDeclaration{LeadingComment: p.consumeComment()}
	decl.Name = p.parseIdentifier()
	p.eat(token.Assign)
	decl.Type = p.parseType()
	return decl
}

func (p *parser) parseEnumDeclaration() *ast.EnumDeclaration {
	p.eat(token.Ident)
	decl := &ast.EnumDeclaration{LeadingComment: p.consumeComment()}
	decl.Name = p.parseIdentifier()
	p.eat(token.LBrace)
	for {
		if p.tok.Kind == token.RBrace {
			break
		}
		decl.Members = append(decl.Members, p.parseEnumMember())
		switch p.tok.Kind {
		case token.Comma:
			p.advance()
		case token.RBrace:
		default:
			panic(fmt.Sprintf("unexpected token %s", p.tok))
		}
	}
	p.eat(token.RBrace)
	return decl
}

func (p *parser) parseInterfaceDeclaration() *ast.InterfaceDeclaration {
	p.eat(token.Ident)
	decl := &ast.InterfaceDeclaration{LeadingComment: p.consumeComment()}
	decl.Name = p.parseIdentifier()
	if p.tok.Kind == token.LAngle {
		decl.TypeParameters = p.parseTypeParameters()
	}
	if p.tok.Kind == token.Ident && p.tok.Text == "extends" {
		decl.HeritageClauses = p.parseHeritageClauses()
	}
	p.eat(token.LBrace)
	for {
		if p.tok.Kind == token.RBrace {
			break
		}
		decl.Members = append(decl.Members, p.parseSignature())
	}
	p.eat(token.RBrace)
	return decl
}

func (p *parser) parseSignature() ast.Signature {
	switch p.tok.Kind {
	case token.Ident:
		return p.parsePropertySignature()
	case token.LBrack:
		return p.parseIndexSignature()
	default:
		panic(fmt.Sprintf("unexpected token %s", p.tok))
	}
}

func (p *parser) parseHeritageClauses() []*ast.HeritageClause {
	p.eat(token.Ident)
	var heritageClauses []*ast.HeritageClause
	for {
		heritageClauses = append(heritageClauses, p.parseHeritageClause())
		if p.tok.Kind != token.Comma {
			break
		}
		p.advance()
	}
	return heritageClauses
}

func (p *parser) parseHeritageClause() *ast.HeritageClause {
	return &ast.HeritageClause{
		Types: []*ast.ExpressionWithTypeArguments{p.parseExpressionWithTypeArguments()},
	}
}

func (p *parser) parseExpressionWithTypeArguments() *ast.ExpressionWithTypeArguments {
	return &ast.ExpressionWithTypeArguments{Expression: p.parseIdentifier()}
}

func (p *parser) parseTypeParameters() []*ast.TypeParameter {
	p.eat(token.LAngle)
	var typeParameters []*ast.TypeParameter
	for {
		typeParameters = append(typeParameters, p.parseTypeParameter())
		if p.tok.Kind != token.Comma {
			break
		}
		p.advance()
	}
	p.eat(token.RAngle)
	return typeParameters
}

func (p *parser) parseTypeParameter() *ast.TypeParameter {
	return &ast.TypeParameter{Name: p.parseIdentifier()}
}

func (p *parser) parsePropertySignature() *ast.PropertySignature {
	signature := &ast.PropertySignature{LeadingComment: p.consumeComment()}
	if p.tok.Kind == token.Ident && p.tok.Text == "readonly" {
		p.advance()
	}
	signature.Name = p.parseIdentifier()
	if p.tok.Kind == token.Question {
		signature.QuestionToken = true
		p.advance()
	}
	p.eat(token.Colon)
	signature.Type = p.parseType()
	if p.tok.Kind == token.Semicolon {
		p.advance()
	}
	signature.TrailingComment = p.consumeLineComment()
	return signature
}

func (p *parser) parseEnumMember() *ast.EnumMember {
	member := &ast.EnumMember{
		Name:           p.parseIdentifier(),
		LeadingComment: p.consumeComment(),
	}
	if p.tok.Kind != token.Assign {
		return member
	}
	p.advance()
	member.Initializer = p.parseInitializer()
	return member
}

func (p *parser) parseInitializer() ast.Expr {
	switch p.tok.Kind {
	case token.Number:
		return &ast.NumericLiteral{Text: p.eat(token.Number).Text}
	case token.String:
		return &ast.StringLiteral{Text: p.eat(token.String).Text}
	case token.Minus:
		p.advance()
		return &ast.PrefixUnaryExpression{
			Operator: token.Minus,
			Operand:  p.parseInitializer(),
		}
	case token.Ident:
		return &ast.TypeReference{TypeName: p.parseIdentifier()}
	case token.LBrack:
		return p.parseArrayLiteralExpression()
	default:
		panic(fmt.Sprintf("unexpected token %s", p.tok))
	}
}

func (p *parser) parseArrayLiteralExpression() *ast.ArrayLiteralExpression {
	p.eat(token.LBrack)
	expr := &ast.ArrayLiteralExpression{}
	for p.tok.Kind != token.RBrack {
		expr.Elements = append(expr.Elements, p.parseInitializer())
		if p.tok.Kind == token.Comma {
			p.advance()
		}
	}
	p.eat(token.RBrack)
	return expr
}

func (p *parser) parseIdentifier() *ast.Identifier {
	return &ast.Identifier{Text: p.eat(token.Ident).Text}
}

func (p *parser) parseType() ast.Type {
	return p.parseTypeCheckUnion()
}

func (p *parser) parseTypeCheckUnion() ast.Type {
	typ := p.parseTypeCheckArray()
	if p.tok.Kind != token.Or {
		return typ
	}

	types := []ast.Type{typ}
	for {
		p.eat(token.Or)
		types = append(types, p.parseTypeCheckArray())
		if p.tok.Kind != token.Or {
			return &ast.UnionType{Types: types}
		}
	}
}

func (p *parser) parseTypeCheckArray() ast.Type {
	typ := p.parseTypeInner()
	if p.tok.Kind != token.LBrack {
		return typ
	}

	p.eat(token.LBrack)
	p.eat(token.RBrack)
	return &ast.ArrayType{ElementType: typ}
}

func (p *parser) parseTypeInner() ast.Type {
	switch p.tok.Kind {
	case token.Ident:
		return p.parseTypeReference()
	case token.LBrace:
		return p.parseTypeLiteral()
	case token.LParen:
		return p.parseParenthesizedType()
	case token.LBrack:
		return p.parseTupleType()
	case token.String:
		return &ast.LiteralType{Literal: &ast.StringLiteral{Text: p.eat(token.String).Text}}
	case token.Number:
		return &ast.LiteralType{Literal: &ast.NumericLiteral{Text: p.eat(token.Number).Text}}
	default:
		panic(fmt.Sprintf("unexpected token %s", p.tok))
	}
}

func (p *parser) parseTupleType() *ast.TupleType {
	p.eat(token.LBrack)
	els := []ast.Type{}
	for {
		els = append(els, p.parseType())
		if p.tok.Kind != token.Comma {
			p.eat(token.RBrack)
			return &ast.TupleType{Elements: els}
		}
		p.advance()
	}
}

func (p *parser) parseTypeReference() *ast.TypeReference {
	first := p.parseIdentifier()
	if p.tok.Kind != token.Dot {
		return &ast.TypeReference{TypeName: first}
	}
	p.advance()
	return &ast.TypeReference{TypeName: &ast.QualifiedName{Left: first, Right: p.parseIdentifier()}}
}

func (p *parser) parseParenthesizedType() *ast.ParenthesizedType {
	p.eat(token.LParen)
	typ := p.parseType()
	p.eat(token.RParen)
	return &ast.ParenthesizedType{Type: typ}
}

func (p *parser) parseTypeLiteral() *ast.TypeLiteral {
	p.eat(token.LBrace)
	literal := &ast.TypeLiteral{}
	for {
		if p.tok.Kind == token.RBrace {
			p.eat(token.RBrace)
			return literal
		}
		literal.Members = append(literal.Members, p.parseSignature())
	}
}

func (p *parser) parseIndexSignature() *ast.IndexSignature {
	signature := &ast.IndexSignature{LeadingComment: p.consumeComment()}
	p.eat(token.LBrack)
	for p.tok.Kind != token.RBrack {
		signature.Parameters = append(signature.Parameters, p.parseParameter())
	}
	p.consumeLineComment() // Throw away
	p.eat(token.RBrack)
	p.eat(token.Colon)
	signature.Type = p.parseType()
	if p.tok.Kind == token.Semicolon {
		p.advance()
	}
	return signature
}

func (p *parser) parseParameter() *ast.Parameter {
	name := p.parseIdentifier()
	p.eat(token.Colon)
	typ := p.parseType()
	return &ast.Parameter{Name: name, Type: typ}
}

func (p *parser) eat(kind token.Kind) token.Token {
	p.expect(kind)
	tok := p.tok
	p.advance()
	return tok
}

func (p *parser) expect(kind token.Kind) {
	if p.tok.Kind != kind {
		panic(fmt.Sprintf("expected kind %s, got %s", kind, p.tok))
	}
}

func (p *parser) advance() {
	for {
		p.tok = p.lex.Pop()
		switch p.tok.Kind {
		case token.Comment:
			p.lastComment = p.tok.Text
		case token.LineComment:
			p.lastLineComment = p.tok.Text
		default:
			return
		}
	}
}

func (p *parser) consumeComment() string {
	comment := p.lastComment
	p.lastComment = ""
	return comment
}

func (p *parser) consumeLineComment() string {
	comment := p.lastLineComment
	p.lastLineComment = ""
	return comment
}
