package main

import (
	"bytes"

	"github.com/agustin-del-pino/aldana/pkg/aldana"
	"github.com/agustin-del-pino/aldana/pkg/aldana/lexer"
	"github.com/agustin-del-pino/aldana/pkg/aldana/ranges"
)

type TokenType int

const (
	Num TokenType = iota
	Eql
	Comma
	LeftBrace
	RightBrace
	LeftPrt
	RightPrt
	Word
	Str
	Hash
)

var (
	NumRange      = ranges.ByteBounded(0x30, 0x39)
	WordRange     = ranges.RangeByteOfRange(ranges.ByteBounded(0x61, 0x7A), ranges.ByteBounded(0x41, 0x5A), ranges.ByteSingle(0x5F))
	AlphaNumRange = ranges.RangeByteOfRange(WordRange, NumRange)
	SpecialRange  = ranges.ByteSet('=', ',', '(', ')', '{', '}', '#')
	QuoteRange    = ranges.ByteSingle(0x22)

	SpecialTokenType = map[byte]TokenType{
		'=': Eql,
		',': Comma,
		'(': LeftPrt,
		')': RightPrt,
		'{': LeftBrace,
		'}': RightBrace,
		'#': Hash,
	}
)

type Token struct {
	Type   TokenType
	Value  []byte
	Line   int
	Column int
}

func IsTokenType(t *Token, p TokenType) bool {
	return t.Type == p
}

func IsTokenValue(t *Token, s string) bool {
	return bytes.Equal(t.Value, []byte(s))
}

func lexNumbs(c lexer.Cursor, r ranges.ByteRange) *Token {
	ln, cl := c.GetPosition()

	t := &Token{
		Type:   Num,
		Line:   ln,
		Column: cl,
	}

	for c.HasChar() && r(c.GetChar()) {
		t.Value = append(t.Value, c.GetChar())
		c.Next()
	}

	return t
}

func lexWord(c lexer.Cursor, _ ranges.ByteRange) *Token {
	ln, cl := c.GetPosition()

	t := &Token{
		Type:   Word,
		Line:   ln,
		Column: cl,
	}

	for c.HasChar() && AlphaNumRange(c.GetChar()) {
		t.Value = append(t.Value, c.GetChar())
		c.Next()
	}

	return t
}

func lexSpecial(c lexer.Cursor, _ ranges.ByteRange) *Token {
	ln, cl := c.GetPosition()

	t := &Token{
		Line:   ln,
		Column: cl,
		Type:   SpecialTokenType[c.GetChar()],
	}

	c.Next()
	return t
}

func lexStr(c lexer.Cursor, r ranges.ByteRange) *Token {
	ln, cl := c.GetPosition()

	t := &Token{
		Type:   Str,
		Line:   ln,
		Column: cl,
	}

	for c.HasChar() && !r(c.GetChar()) {
		t.Value = append(t.Value, c.GetChar())
		c.Next()
	}

	return t
}

func NewPythonLexer() lexer.Lexer[*Token] {
	return aldana.NewLexer(&aldana.LexerOptions[*Token]{
		Ignore: aldana.IgnoreWhiteSpaces(),
		LexRules: []aldana.LexicalRule[*Token]{
			aldana.NewLexicalRule(NumRange, lexNumbs),
			aldana.NewLexicalRule(WordRange, lexWord),
			aldana.NewLexicalRule(SpecialRange, lexSpecial),
			aldana.NewLexicalRule(QuoteRange, lexStr),
		},
	})
}
