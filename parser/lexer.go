package parser

import (
	"bytes"

	"github.com/armsnyder/typescript-ast-go/token"
)

type lexer struct {
	Source []byte

	offset                int
	isInsideBlock         bool
	willBeTrailingComment bool
	nextToken             token.Token
}

func (x *lexer) Peek() token.Token {
	if x.nextToken.Kind == 0 {
		x.nextToken = x.next()
	}

	return x.nextToken
}

func (x *lexer) Pop() token.Token {
	if x.nextToken.Kind == 0 {
		return x.next()
	}

	t := x.nextToken
	x.nextToken = token.Token{}
	return t
}

func (x *lexer) next() token.Token {
	for x.offset < len(x.Source) {
		switch x.Source[x.offset] {
		case ' ', '\t', '\r':
			x.offset++

		case '\n':
			x.offset++
			x.willBeTrailingComment = false

		case '/':
			return x.nextComment()

		default:
			x.willBeTrailingComment = true

			switch x.Source[x.offset] {
			case '|':
				return x.char(token.Or)

			case '=':
				return x.char(token.Assign)

			case '-':
				return x.char(token.Minus)

			case '(':
				return x.char(token.LParen)

			case ')':
				return x.char(token.RParen)

			case '[':
				return x.char(token.LBrack)

			case ']':
				return x.char(token.RBrack)

			case '{':
				x.isInsideBlock = true
				return x.char(token.LBrace)

			case '}':
				x.isInsideBlock = false
				return x.char(token.RBrace)

			case '<':
				return x.char(token.LAngle)

			case '>':
				return x.char(token.RAngle)

			case ',':
				return x.char(token.Comma)

			case '.':
				return x.char(token.Dot)

			case ':':
				return x.char(token.Colon)

			case ';':
				return x.char(token.Semicolon)

			case '?':
				return x.char(token.Question)

			case '\'':
				return x.nextString()

			default:
				if x.Source[x.offset] >= '0' && x.Source[x.offset] <= '9' {
					return x.nextNumber()
				}

				if (x.Source[x.offset] >= 'a' && x.Source[x.offset] <= 'z') || (x.Source[x.offset] >= 'A' && x.Source[x.offset] <= 'Z') {
					return x.nextIdent()
				}

				return token.Token{Kind: token.Illegal}
			}
		}
	}

	return token.Token{Kind: token.EOF}
}

func (x *lexer) nextComment() token.Token {
	if x.offset+1 >= len(x.Source) {
		return token.Token{Kind: token.Illegal}
	}

	switch x.Source[x.offset+1] {
	case '/':
		return x.nextLineComment()

	case '*':
		return x.nextBlockComment()

	default:
		return token.Token{Kind: token.Illegal}
	}
}

func (x *lexer) nextLineComment() token.Token {
	x.offset += 2
	for x.offset+1 < len(x.Source) && x.Source[x.offset] == ' ' {
		x.offset++
	}

	commentStart := x.offset

	for x.offset < len(x.Source) && x.Source[x.offset] != '\n' {
		x.offset++
	}

	kind := token.Comment
	if x.willBeTrailingComment {
		kind = token.LineComment
	}

	return token.Token{Kind: kind, Text: string(x.Source[commentStart:x.offset])}
}

func (x *lexer) nextBlockComment() token.Token {
	x.offset += 2
	for x.offset+1 < len(x.Source) && x.Source[x.offset] == '*' {
		x.offset++
	}

	innerEndIndex := bytes.Index(x.Source[x.offset:], []byte("*/"))
	if innerEndIndex == -1 {
		return token.Token{Kind: token.Illegal}
	}
	innerEndIndex += x.offset
	endIndex := innerEndIndex + 2
	for innerEndIndex > 0 && x.Source[innerEndIndex-1] == '*' {
		innerEndIndex--
	}

	var comment []byte

	for x.offset < innerEndIndex {
		lineEnd := bytes.Index(x.Source[x.offset:innerEndIndex], []byte("\n"))
		if lineEnd == -1 {
			comment = append(comment, x.Source[x.offset:innerEndIndex]...)
			break
		}

		comment = append(comment, x.Source[x.offset:x.offset+lineEnd]...)
		comment = append(comment, '\n')
		x.offset += lineEnd + 1

		i := bytes.IndexFunc(x.Source[x.offset:innerEndIndex], func(r rune) bool {
			return r != ' ' && r != '\t' && r != '*'
		})
		if i == -1 {
			break
		}

		x.offset += i
	}

	x.offset = endIndex

	kind := token.Comment
	if x.willBeTrailingComment {
		kind = token.LineComment
	}

	return token.Token{Kind: kind, Text: string(bytes.TrimSpace(comment))}
}

func (x *lexer) nextString() token.Token {
	x.offset++
	start := x.offset

	end := bytes.IndexByte(x.Source[x.offset:], '\'')
	if end == -1 {
		return token.Token{Kind: token.Illegal}
	}

	x.offset += end + 1

	return token.Token{Kind: token.String, Text: string(x.Source[start : x.offset-1])}
}

func (x *lexer) char(kind token.Kind) token.Token {
	x.offset++
	return token.Token{Kind: kind}
}

func (x *lexer) nextNumber() token.Token {
	start := x.offset
	for x.offset < len(x.Source) && x.Source[x.offset] >= '0' && x.Source[x.offset] <= '9' {
		x.offset++
	}
	return token.Token{Kind: token.Number, Text: string(x.Source[start:x.offset])}
}

func (x *lexer) nextIdent() (tok token.Token) {
	start := x.offset

	end := bytes.IndexFunc(x.Source[x.offset:], func(r rune) bool {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_':
			return false
		default:
			return true
		}
	})
	if end == -1 {
		x.offset = len(x.Source)
	} else {
		x.offset += end
	}

	value := string(x.Source[start:x.offset])

	return token.Token{Kind: token.Ident, Text: value}
}
