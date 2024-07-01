// Package token contains the definition of the lexical tokens of the
// TypeScript programming language.
package token

import "strconv"

// Kind is the set of lexical tokens of the TypeScript programming
// language.
type Kind int

const (
	// Special tokens.
	Illegal Kind = iota
	EOF

	// Comments.
	Comment     // // or /* */ at the beginning of a line
	LineComment // // or /* */ at the end of a line

	// Identifiers and literals.
	Ident  // main, const, extends, etc.
	Number // 12345
	String // "abc"

	// Operators.
	Or     // |
	Assign // =
	Minus  // -

	// Delimiters and punctuation.
	LParen    // (
	RParen    // )
	LBrack    // [
	RBrack    // ]
	LBrace    // {
	RBrace    // }
	LAngle    // <
	RAngle    // >
	Comma     // ,
	Dot       // .
	Colon     // :
	Semicolon // ;
	Question  // ?
)

var tokens = [...]string{
	// Special tokens.
	Illegal: "Illegal",
	EOF:     "EOF",

	// Comments.
	Comment:     "Comment",
	LineComment: "LineComment",

	// Identifiers and literals.
	Ident:  "Ident",
	Number: "Number",
	String: "String",

	// Operators.
	Or:     "|",
	Assign: "=",
	Minus:  "-",

	// Delimiters and punctuation.
	LParen:    "(",
	RParen:    ")",
	LBrack:    "[",
	RBrack:    "]",
	LBrace:    "{",
	RBrace:    "}",
	LAngle:    "<",
	RAngle:    ">",
	Comma:     ",",
	Dot:       ".",
	Colon:     ":",
	Semicolon: ";",
	Question:  "?",
}

func (k Kind) String() string {
	if 0 <= k && k < Kind(len(tokens)) {
		return tokens[k]
	}
	return "token(" + strconv.Itoa(int(k)) + ")"
}

func (k Kind) IsLiteral() bool {
	switch k {
	case Ident, Number, String, Comment, LineComment:
		return true
	default:
		return false
	}
}

// Token represents a lexical token of the TypeScript programming language.
// Identifiers and literals are accompanied by the corresponding text.
type Token struct {
	Kind Kind
	Text string
}
