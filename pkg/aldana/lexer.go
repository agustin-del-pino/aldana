package aldana

import (
	"github.com/agustin-del-pino/aldana/pkg/aldana/lexer"
	"github.com/agustin-del-pino/aldana/pkg/aldana/ranges"
)

// TokenRule is a function that returns a new token based on some defined rules.
type TokenRule[T any] func(c lexer.Cursor, r ranges.ByteRange) T

// LexicalRule is a function that returns a ranges.ByteRange related to a TokenRule.
type LexicalRule[T any] func() (ranges.ByteRange, TokenRule[T])

// LexicalOmit is a function that returns a range.ByteRange to omit and the way of how the bytes have to be omitted.
type LexicalOmit func() (ranges.ByteRange, func(c lexer.Cursor, r ranges.ByteRange))

// LexerOptions contains the options for configure the default implementation of the Lexer.
type LexerOptions[T any] struct {
	// LexRules is a ordered slice of LexicalRule.
	LexRules []LexicalRule[T]
	// Ignore is a pseudo-lexical-rule that ignores chars.
	Ignore LexicalOmit
}

// defaultLexer implements lexer.Lexer
type defaultLexer[T any] struct {
	ops *LexerOptions[T]
}

func (l *defaultLexer[T]) Tokenize(c lexer.Cursor) ([]T, error) {
	tks := []T{}

	c.Next()

	for c.HasChar() {
		if r, i := l.ops.Ignore(); r(c.GetChar()) {
			i(c, r)
			continue
		}
		var ux bool

		for _, lr := range l.ops.LexRules {
			if r, t := lr(); r(c.GetChar()) {
				tks = append(tks, t(c, r))
				ux = false
				break
			}
			ux = true
		}

		if ux {
			return nil, lexer.ErrUnexpectedChar
		}
	}

	return tks, nil
}

// NewLexer returns the default implementation of lexer.Lexer.
//
// # About the implementation
//   - The priority is: ignore then lex-rules. And those rules are a ordered slice of LexicalRule.
//   - When the character is not consumed by any lex-rule or ignored, the lexer.ErrUnexpectedChar is returned.
//
// # Example
//
//	type Token struct {
//		Type string
//		Value []byte
//	}
//
//	func lexNumbs(c lexer.Cursor, r ranges.ByteRange) *Token {
//		t := &Token{Type: "num"}
//
//		for c.HasChar() && r(c.GetChar()) {
//			t.Value = append(t.Value, c.GetChar())
//			c.Next()
//		}
//
//		return t
//	}
//
//	l := NewLexer[*Token](&LexerOptions{
//		Ignore: IgnoreWhiteSpaces(),
//		LexRules: []LexicalRules {
//			NewLexicalRule(range.BytesBounded(0x30, 0x39), lexNumbs)
//		}
//	})
//
//	tks := l.Tokenize(NewCursor([]byte("123456789 4450048 777")))
func NewLexer[T any](ops *LexerOptions[T]) lexer.Lexer[T] {
	return &defaultLexer[T]{
		ops: ops,
	}
}

// NewLexicalRule returns a LexicalRule. Use as short-cut.
func NewLexicalRule[T any](r ranges.ByteRange, t TokenRule[T]) LexicalRule[T] {
	return func() (ranges.ByteRange, TokenRule[T]) {
		return r, t
	}
}

// IgnoreWhiteSpaces returns a LexicalOmit for ignore the 0x20 (white-space).
func IgnoreWhiteSpaces() LexicalOmit {
	return func() (ranges.ByteRange, func(c lexer.Cursor, r ranges.ByteRange)) {
		return ranges.ByteSingle(0x20), func(c lexer.Cursor, r ranges.ByteRange) {
			c.Next()
		}
	}
}
